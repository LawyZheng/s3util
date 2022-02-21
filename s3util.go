package s3util

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/lawyzheng/s3util/uploader"
)

type S3Client struct {
	config   *S3Config
	client   *s3.S3
	updriver *s3manager.Uploader
	uploader uploader.Uploader
}

type S3Config struct {
	forcePath     bool
	endpoint      string
	accessKey     string
	secretKey     string
	token         string
	region        string
	uploadTimeout time.Duration
}

func NewDefaultConfig(endpoint, ak, sk string) *S3Config {
	return &S3Config{
		forcePath: true,
		endpoint:  endpoint,
		accessKey: ak,
		secretKey: sk,
		region:    "us-east-1",
	}
}

func (c *S3Config) SetUploadTimeOut(t time.Duration) {
	c.uploadTimeout = t
}

func (s3client *S3Client) initClient() {
	sess := session.Must(session.NewSession())
	cfg := aws.NewConfig()
	cfg.WithEndpoint(s3client.config.endpoint)
	cfg.WithRegion(s3client.config.region)
	cfg.WithS3ForcePathStyle(s3client.config.forcePath)
	creds := credentials.NewStaticCredentials(
		s3client.config.accessKey,
		s3client.config.secretKey,
		s3client.config.token,
	)
	cfg.WithCredentials(creds)
	svc := s3.New(sess, cfg)
	s3client.client = svc
	s3client.updriver = s3manager.NewUploaderWithClient(svc)
}

func NewS3Client(endpoint, ak, sk string) *S3Client {
	config := NewDefaultConfig(endpoint, ak, sk)
	return NewS3ClientWithConfig(config)
}

func NewS3ClientWithConfig(config *S3Config) *S3Client {
	client := &S3Client{
		config: config,
	}
	client.initClient()
	return client
}

func (c *S3Client) getUploadDriver() *s3manager.Uploader {
	return c.updriver
}

func (c *S3Client) GetUploadDriver() *s3manager.Uploader {
	return c.updriver
}

func (c *S3Client) SetUploadTimeOut(t time.Duration) {
	c.config.SetUploadTimeOut(t)
}

func (c *S3Client) CountObjectWithPrefix(bucketName, prefix string) (int, error) {
	if !strings.HasSuffix(prefix, "/") {
		prefix = prefix + "/"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result := 0
	err := c.client.ListObjectsV2PagesWithContext(ctx, &s3.ListObjectsV2Input{
		Bucket: &bucketName,
		Prefix: &prefix,
	}, func(p *s3.ListObjectsV2Output, last bool) (shouldContinue bool) {
		for _, obj := range p.Contents {
			if *obj.Key != prefix {
				result++
			}
		}
		return true
	})
	return result, err
}

func (c *S3Client) GetObjectDownloadURL(bucketName string, key string, expire time.Duration) (string, error) {
	params := &s3.GetObjectInput{
		Bucket: &bucketName,
		Key:    &key,
	}
	req, _ := c.client.GetObjectRequest(params)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req.SetContext(ctx)
	url, err := req.Presign(expire)
	return url, err
}

func (c *S3Client) GetHeadObject(objKey, bucket string) (*s3.HeadObjectOutput, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.client.HeadObjectWithContext(
		ctx,
		&s3.HeadObjectInput{
			Bucket: &bucket,
			Key:    &objKey,
		})
}

func (c *S3Client) SetUploader(up uploader.Uploader) {
	if up == nil {
		panic("uploader is <nil>")
	}
	c.uploader = up
}

func (c *S3Client) UploadHttpResponse(bucket, objKey string, resp *http.Response) (bool, error) {
	var upload *uploader.HttpUploader
	if c.uploader != nil {
		if up, ok := c.uploader.(*uploader.HttpUploader); ok {
			upload = up
		}
	}

	if upload == nil {
		upload = uploader.NewSimpleHttpUploader(c.getUploadDriver())
		upload.SetTimeout(c.config.uploadTimeout)
	}

	head, err := c.GetHeadObject(objKey, bucket)
	if err != nil {
		return false, fmt.Errorf("GetHeadObjectFailed: caused by: %s", err)
	}

	tag := resp.Header.Get("ETag")
	exist := upload.CheckObjectExist(tag, head)
	if exist {
		return true, nil
	}

	return false, upload.UploadObject(bucket, objKey, resp.Body, tag)
}

func (c *S3Client) CheckObjectExist(bucket, objKey, srcEtag string) bool {
	if c.uploader == nil {
		panic("uploader is <nil>")
	}

	if strings.TrimSpace(srcEtag) == "" {
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	resp, err := c.client.HeadObjectWithContext(
		ctx,
		&s3.HeadObjectInput{
			Bucket: &bucket,
			Key:    &objKey,
		})
	if err != nil {
		return false
	}

	return c.uploader.CheckObjectExist(srcEtag, resp)
}

func (c *S3Client) UploadObject(bucketName, key string, src io.Reader, tag string) error {
	if c.uploader == nil {
		panic("uploader is <nil>")
	}

	return c.uploader.UploadObject(bucketName, key, src, tag)
}

func (c *S3Client) UploadObjetWithMetadata(bucketName, key string, src io.Reader, meta map[string]*string) error {
	if c.uploader == nil {
		panic("uploader is <nil>")
	}

	return c.uploader.UploadObjetWithMetadata(bucketName, key, src, meta)
}

func (c *S3Client) PutObjectWithMetadata(bucketName, key string, src io.ReadSeeker, meta map[string]*string) error {
	params := &s3.PutObjectInput{
		Bucket:   aws.String(bucketName),
		Key:      aws.String(key),
		Body:     src,
		Metadata: meta,
	}

	var ctx context.Context
	var cancel context.CancelFunc
	if c.config.uploadTimeout != 0 {
		ctx, cancel = context.WithTimeout(context.Background(), c.config.uploadTimeout)
		defer cancel()
	} else {
		ctx = context.Background()
	}

	_, err := c.client.PutObjectWithContext(ctx, params)
	return err
}

func (c *S3Client) CreateBucket(bucketName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := c.client.CreateBucketWithContext(
		ctx,
		&s3.CreateBucketInput{
			Bucket: &bucketName,
		},
	)
	if aerr, ok := err.(awserr.Error); ok {
		switch aerr.Code() {
		case s3.ErrCodeBucketAlreadyOwnedByYou, s3.ErrCodeBucketAlreadyExists:
			return nil
		}
	}
	return err
}

func (c *S3Client) BucketExist(bucketName string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := c.client.HeadBucketWithContext(
		ctx,
		&s3.HeadBucketInput{
			Bucket: &bucketName,
		},
	)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchBucket, "NotFound":
				return false, nil
			}
		}
	}

	return err == nil, err
}
