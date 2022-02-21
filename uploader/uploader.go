package uploader

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

const (
	EtagCheck = "S3utiletagchk"
)

type CheckObjectExistFn func(local string, remote *s3.HeadObjectOutput) bool

type Uploader interface {
	CheckObjectExist(local string, remote *s3.HeadObjectOutput) bool
	UploadObject(bucket, key string, source io.Reader, tag string) error
	UploadObjetWithMetadata(bucket, key string, source io.Reader, meta map[string]*string) error
}

func parseEtagChk(s string) string {
	return strings.Title(strings.ToLower(s))
}

func NewChkFnWithEtagChk(etagchk string) CheckObjectExistFn {
	etagchk = strings.TrimSpace(etagchk)
	if etagchk == "" {
		panic("etagchk can't be empty")
	}
	etagchk = parseEtagChk(etagchk)

	return func(local string, remote *s3.HeadObjectOutput) bool {
		s := strings.TrimSpace(local)
		if s == "" {
			return false
		}

		return s ==
			strings.TrimSpace(aws.StringValue(remote.Metadata[etagchk]))
	}
}

type HttpUploader struct {
	S3uploader *s3manager.Uploader
	CheckFn    CheckObjectExistFn
	Timeout    time.Duration
	EtagCheck  string
}

func NewSimpleHttpUploader(uploader *s3manager.Uploader) *HttpUploader {
	return &HttpUploader{
		S3uploader: uploader,
	}
}

func (h *HttpUploader) getEtagChk() string {
	if s := strings.TrimSpace(h.EtagCheck); s != "" {
		return parseEtagChk(s)
	}
	return EtagCheck
}

func (h *HttpUploader) getChkFn() CheckObjectExistFn {
	if h.CheckFn != nil {
		return h.CheckFn
	}
	return NewChkFnWithEtagChk(h.getEtagChk())
}

func (h *HttpUploader) SetUploader(up *s3manager.Uploader) {
	if up == nil {
		panic("s3 uploader is <nil>")
	}
	h.S3uploader = up
}

func (h *HttpUploader) SetTimeout(t time.Duration) {
	h.Timeout = t
}

func (h *HttpUploader) SetEtagCheck(etagchk string) {
	s := strings.TrimSpace(etagchk)
	if s == "" {
		panic("etagchk can't be empty")
	}
	h.EtagCheck = s
}

func (h *HttpUploader) CheckObjectExist(local string, remote *s3.HeadObjectOutput) bool {
	return h.getChkFn()(local, remote)
}

func (h *HttpUploader) UploadObject(bucket, key string, source io.Reader, tag string) error {
	t := time.Now().Format("2006-01-02 15:04:05")
	meta := map[string]*string{
		h.getEtagChk(): aws.String(tag),
		"modifiedtime": &t,
	}
	return h.UploadObjetWithMetadata(bucket, key, source, meta)
}

func (h *HttpUploader) UploadObjetWithMetadata(bucket, key string, source io.Reader, meta map[string]*string) error {
	if h.S3uploader == nil {
		panic("s3 uploader is <nil>")
	}

	params := &s3manager.UploadInput{
		Bucket:   aws.String(bucket),
		Key:      aws.String(key),
		Body:     source,
		Metadata: meta,
	}

	var ctx context.Context
	var cancel context.CancelFunc
	if h.Timeout != 0 {
		ctx, cancel = context.WithTimeout(context.Background(), h.Timeout)
		defer cancel()
	} else {
		ctx = context.Background()
	}

	_, err := h.S3uploader.UploadWithContext(ctx, params)
	if err != nil {
		return fmt.Errorf("UploadFailed: caused by: %s", err)
	}
	return nil
}
