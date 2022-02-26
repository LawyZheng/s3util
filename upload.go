package s3util

import (
	"context"
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

const (
	etagCheck = "S3utiletagchk"
)

type HttpChkFn func(src http.Header, remote *s3.HeadObjectOutput) bool

type HttpChkInterface interface {
	GetChkFn() HttpChkFn
	GetEtag() string
}

func parseEtagChk(s string) string {
	return strings.Title(strings.ToLower(s))
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
	etagchk = parseEtagChk(etagchk)

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
