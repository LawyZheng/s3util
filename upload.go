package s3util

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type UploadResponse struct {
	Skip   bool
	Output *s3manager.UploadOutput
}

type UploadRequest struct {
	src  io.Reader
	ctx  context.Context
	meta map[string]*string
}

func NewUploadRequest(src io.Reader, ctx context.Context, meta map[string]*string) *UploadRequest {
	return &UploadRequest{
		src:  src,
		ctx:  ctx,
		meta: meta,
	}
}

func parseChkKey(s string) string {
	return strings.Title(strings.ToLower(s))
}

const (
	etagCheck = "S3utiletagchk"
)

type HttpChkFn func(src http.Header, remote *s3.HeadObjectOutput) bool

type HttpChkInterface interface {
	GetChkFn() HttpChkFn
	GetEtag() string
}

func defaultHttpChk() *HttpChk {
	return NewHttpChkWithEtagChk(etagCheck)
}

type HttpChk struct {
	etag string
	fn   func(src http.Header, remote *s3.HeadObjectOutput) bool
}

func (h *HttpChk) GetChkFn() HttpChkFn {
	return h.fn
}
func (h *HttpChk) GetEtag() string {
	return h.etag
}

func NewHttpChkWithEtagChk(etagchk string) *HttpChk {
	etagchk = strings.TrimSpace(etagchk)
	if etagchk == "" {
		panic("etagchk can't be empty")
	}
	etagchk = parseChkKey(etagchk)

	return &HttpChk{
		etag: etagchk,
		fn: func(src http.Header, remote *s3.HeadObjectOutput) bool {
			s := strings.TrimSpace(src.Get("ETag"))
			if s == "" {
				return false
			}

			return s ==
				strings.TrimSpace(aws.StringValue(remote.Metadata[etagchk]))
		},
	}
}

const (
	md5Check = "S3utilmd5chk"
)

type FileChkFn func(src string, remote *s3.HeadObjectOutput) bool

type FileChkInterface interface {
	GetChkFn() FileChkFn
	GetMd5Key() string
}

type FileChk struct {
	md5Key string
	fn     func(src string, remote *s3.HeadObjectOutput) bool
}

func (f *FileChk) GetChkFn() FileChkFn {
	return f.fn
}

func (f *FileChk) GetMd5Key() string {
	return f.md5Key
}

func NewFileChkWithMd5Chk(md5chk string) *FileChk {
	md5chk = strings.TrimSpace(md5chk)
	if md5chk == "" {
		panic("md5chk can't be empty")
	}
	md5chk = parseChkKey(md5chk)

	return &FileChk{
		md5Key: md5chk,
		fn: func(src string, remote *s3.HeadObjectOutput) bool {
			srcMd5 := strings.TrimSpace(src)
			metaChk := strings.TrimSpace(aws.StringValue(remote.Metadata[md5chk]))
			if srcMd5 == metaChk {
				return true
			}

			return aws.StringValue(remote.ETag) ==
				fmt.Sprintf("\"%s\"", srcMd5)
		},
	}
}

func defaultFileChk() *FileChk {
	return NewFileChkWithMd5Chk(md5Check)
}
