package s3uploader

import (
	"fmt"
	"io"
	"reflect"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/stretchr/testify/assert"
)

func getDrvier() *s3manager.Uploader {
	sess := session.Must(session.NewSession())
	cfg := aws.NewConfig()
	cfg.WithEndpoint("")
	cfg.WithRegion("")
	cfg.WithS3ForcePathStyle(true)
	creds := credentials.NewStaticCredentials(
		"",
		"",
		"",
	)
	cfg.WithCredentials(creds)
	svc := s3.New(sess, cfg)
	return s3manager.NewUploaderWithClient(svc)
}

func Test_parseEtagChk(t *testing.T) {
	assert := assert.New(t)

	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			"lower case",
			args{s: "test"},
			"Test",
		},
		{
			"capital case",
			args{s: "Test"},
			"Test",
		},
		{
			"upper case",
			args{s: "TEST"},
			"Test",
		},
		{
			"mix case",
			args{s: "TeSt"},
			"Test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseEtagChk(tt.args.s)
			assert.Equal(tt.want, got, fmt.Sprintln("Case: ", tt.name))
		})
	}
}

func TestNewChkFnWithEtagChk(t *testing.T) {
	type args struct {
		etagchk string
	}
	tests := []struct {
		name string
		args args
		want CheckObjectExistFn
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewChkFnWithEtagChk(tt.args.etagchk); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewChkFnWithEtagChk() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewSimpleHttpUploader(t *testing.T) {
	type args struct {
		uploader *s3manager.Uploader
	}
	tests := []struct {
		name string
		args args
		want *HttpUploader
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSimpleHttpUploader(tt.args.uploader); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSimpleHttpUploader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHttpUploader_GetEtagChk(t *testing.T) {
	assert := assert.New(t)

	type fields struct {
		driver    *s3manager.Uploader
		checkFn   CheckObjectExistFn
		timeout   time.Duration
		etagCheck string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
		{
			"lower case",
			fields{
				driver:    getDrvier(),
				etagCheck: "test",
			},
			"Test",
		},
		{
			"upper case",
			fields{
				driver:    getDrvier(),
				etagCheck: "TEST",
			},
			"Test",
		},
		{
			"capital case",
			fields{
				driver:    getDrvier(),
				etagCheck: "Test",
			},
			"Test",
		},
		{
			"mixed case",
			fields{
				driver:    getDrvier(),
				etagCheck: "TeSt",
			},
			"Test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &HttpUploader{
				driver:    tt.fields.driver,
				checkFn:   tt.fields.checkFn,
				timeout:   tt.fields.timeout,
				etagCheck: tt.fields.etagCheck,
			}

			got := h.GetEtagChk()
			assert.Equal(tt.want, got, fmt.Sprintln("Case: ", tt.name))
		})
	}
}

func TestHttpUploader_SetEtagCheck(t *testing.T) {
	type fields struct {
		driver    *s3manager.Uploader
		checkFn   CheckObjectExistFn
		timeout   time.Duration
		etagCheck string
	}
	type args struct {
		etagchk string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &HttpUploader{
				driver:    tt.fields.driver,
				checkFn:   tt.fields.checkFn,
				timeout:   tt.fields.timeout,
				etagCheck: tt.fields.etagCheck,
			}
			h.SetEtagCheck(tt.args.etagchk)
		})
	}
}

func TestHttpUploader_GetChkFn(t *testing.T) {
	type fields struct {
		driver    *s3manager.Uploader
		checkFn   CheckObjectExistFn
		timeout   time.Duration
		etagCheck string
	}
	tests := []struct {
		name   string
		fields fields
		want   CheckObjectExistFn
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &HttpUploader{
				driver:    tt.fields.driver,
				checkFn:   tt.fields.checkFn,
				timeout:   tt.fields.timeout,
				etagCheck: tt.fields.etagCheck,
			}
			if got := h.GetChkFn(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HttpUploader.GetChkFn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHttpUploader_SetChkFn(t *testing.T) {
	type fields struct {
		driver    *s3manager.Uploader
		checkFn   CheckObjectExistFn
		timeout   time.Duration
		etagCheck string
	}
	type args struct {
		fn CheckObjectExistFn
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &HttpUploader{
				driver:    tt.fields.driver,
				checkFn:   tt.fields.checkFn,
				timeout:   tt.fields.timeout,
				etagCheck: tt.fields.etagCheck,
			}
			h.SetChkFn(tt.args.fn)
		})
	}
}

func TestHttpUploader_GetDriver(t *testing.T) {
	type fields struct {
		driver    *s3manager.Uploader
		checkFn   CheckObjectExistFn
		timeout   time.Duration
		etagCheck string
	}
	tests := []struct {
		name   string
		fields fields
		want   *s3manager.Uploader
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &HttpUploader{
				driver:    tt.fields.driver,
				checkFn:   tt.fields.checkFn,
				timeout:   tt.fields.timeout,
				etagCheck: tt.fields.etagCheck,
			}
			if got := h.GetDriver(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HttpUploader.GetDriver() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHttpUploader_SetDriver(t *testing.T) {
	type fields struct {
		driver    *s3manager.Uploader
		checkFn   CheckObjectExistFn
		timeout   time.Duration
		etagCheck string
	}
	type args struct {
		up *s3manager.Uploader
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &HttpUploader{
				driver:    tt.fields.driver,
				checkFn:   tt.fields.checkFn,
				timeout:   tt.fields.timeout,
				etagCheck: tt.fields.etagCheck,
			}
			h.SetDriver(tt.args.up)
		})
	}
}

func TestHttpUploader_SetTimeout(t *testing.T) {
	type fields struct {
		driver    *s3manager.Uploader
		checkFn   CheckObjectExistFn
		timeout   time.Duration
		etagCheck string
	}
	type args struct {
		t time.Duration
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &HttpUploader{
				driver:    tt.fields.driver,
				checkFn:   tt.fields.checkFn,
				timeout:   tt.fields.timeout,
				etagCheck: tt.fields.etagCheck,
			}
			h.SetTimeout(tt.args.t)
		})
	}
}

func TestHttpUploader_GetTimeout(t *testing.T) {
	type fields struct {
		driver    *s3manager.Uploader
		checkFn   CheckObjectExistFn
		timeout   time.Duration
		etagCheck string
	}
	tests := []struct {
		name   string
		fields fields
		want   time.Duration
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &HttpUploader{
				driver:    tt.fields.driver,
				checkFn:   tt.fields.checkFn,
				timeout:   tt.fields.timeout,
				etagCheck: tt.fields.etagCheck,
			}
			if got := h.GetTimeout(); got != tt.want {
				t.Errorf("HttpUploader.GetTimeout() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHttpUploader_CheckObjectExist(t *testing.T) {
	type fields struct {
		driver    *s3manager.Uploader
		checkFn   CheckObjectExistFn
		timeout   time.Duration
		etagCheck string
	}
	type args struct {
		local  string
		remote *s3.HeadObjectOutput
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &HttpUploader{
				driver:    tt.fields.driver,
				checkFn:   tt.fields.checkFn,
				timeout:   tt.fields.timeout,
				etagCheck: tt.fields.etagCheck,
			}
			if got := h.CheckObjectExist(tt.args.local, tt.args.remote); got != tt.want {
				t.Errorf("HttpUploader.CheckObjectExist() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHttpUploader_UploadObject(t *testing.T) {
	type fields struct {
		driver    *s3manager.Uploader
		checkFn   CheckObjectExistFn
		timeout   time.Duration
		etagCheck string
	}
	type args struct {
		bucket string
		key    string
		source io.Reader
		tag    string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &HttpUploader{
				driver:    tt.fields.driver,
				checkFn:   tt.fields.checkFn,
				timeout:   tt.fields.timeout,
				etagCheck: tt.fields.etagCheck,
			}
			if err := h.UploadObject(tt.args.bucket, tt.args.key, tt.args.source, tt.args.tag); (err != nil) != tt.wantErr {
				t.Errorf("HttpUploader.UploadObject() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHttpUploader_UploadObjetWithMetadata(t *testing.T) {
	type fields struct {
		driver    *s3manager.Uploader
		checkFn   CheckObjectExistFn
		timeout   time.Duration
		etagCheck string
	}
	type args struct {
		bucket string
		key    string
		source io.Reader
		meta   map[string]*string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &HttpUploader{
				driver:    tt.fields.driver,
				checkFn:   tt.fields.checkFn,
				timeout:   tt.fields.timeout,
				etagCheck: tt.fields.etagCheck,
			}
			if err := h.UploadObjetWithMetadata(tt.args.bucket, tt.args.key, tt.args.source, tt.args.meta); (err != nil) != tt.wantErr {
				t.Errorf("HttpUploader.UploadObjetWithMetadata() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
