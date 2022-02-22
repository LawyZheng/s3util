package s3uploader

import (
	"io"
	"time"

	"github.com/aws/aws-sdk-go/service/s3"
)

type CheckObjectExistFn func(local string, remote *s3.HeadObjectOutput) bool

type Uploader interface {
	SetTimeout(t time.Duration)
	CheckObjectExist(local string, remote *s3.HeadObjectOutput) bool
	UploadObject(bucket, key string, source io.Reader, tag string) error
	UploadObjetWithMetadata(bucket, key string, source io.Reader, meta map[string]*string) error
}
