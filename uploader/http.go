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
	etagCheck = "S3utiletagchk"
)

type HttpUploader struct {
	driver    *s3manager.Uploader
	checkFn   CheckObjectExistFn
	timeout   time.Duration
	etagCheck string
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

func NewSimpleHttpUploader(uploader *s3manager.Uploader) *HttpUploader {
	return &HttpUploader{
		driver: uploader,
	}
}

func (h *HttpUploader) GetEtagChk() string {
	if s := strings.TrimSpace(h.etagCheck); s != "" {
		return parseEtagChk(s)
	}
	return etagCheck
}

func (h *HttpUploader) SetEtagCheck(etagchk string) {
	s := strings.TrimSpace(etagchk)
	if s == "" {
		panic("etagchk can't be empty")
	}
	h.etagCheck = s
}

func (h *HttpUploader) GetChkFn() CheckObjectExistFn {
	if h.checkFn != nil {
		return h.checkFn
	}
	return NewChkFnWithEtagChk(h.GetEtagChk())
}

func (h *HttpUploader) SetChkFn(fn CheckObjectExistFn) {
	if fn == nil {
		panic("CheckObjectExistFn is <nil>")
	}
	h.checkFn = fn
}

func (h *HttpUploader) GetDriver() *s3manager.Uploader {
	return h.driver
}

func (h *HttpUploader) SetDriver(up *s3manager.Uploader) {
	if up == nil {
		panic("s3 uploader driver is <nil>")
	}
	h.driver = up
}

func (h *HttpUploader) SetTimeout(t time.Duration) {
	h.timeout = t
}

func (h *HttpUploader) GetTimeout() time.Duration {
	return h.timeout
}

func (h *HttpUploader) CheckObjectExist(local string, remote *s3.HeadObjectOutput) bool {
	return h.GetChkFn()(local, remote)
}

func (h *HttpUploader) UploadObject(bucket, key string, source io.Reader, tag string) error {
	t := time.Now().Format("2006-01-02 15:04:05")
	meta := map[string]*string{
		h.GetEtagChk(): aws.String(tag),
		"modifiedtime": &t,
	}
	return h.UploadObjetWithMetadata(bucket, key, source, meta)
}

func (h *HttpUploader) UploadObjetWithMetadata(bucket, key string, source io.Reader, meta map[string]*string) error {
	if h.GetDriver() == nil {
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
	if h.GetTimeout() != 0 {
		ctx, cancel = context.WithTimeout(context.Background(), h.GetTimeout())
		defer cancel()
	} else {
		ctx = context.Background()
	}

	_, err := h.GetDriver().UploadWithContext(ctx, params)
	if err != nil {
		return fmt.Errorf("UploadFailed: caused by: %s", err)
	}
	return nil
}
