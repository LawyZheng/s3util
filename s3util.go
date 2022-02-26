package s3util

import (
	"context"
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
)

type S3Client struct {
	config   *S3Config
	client   *s3.S3
	updriver *s3manager.Uploader

	// uploader s3uploader.Uploader
	HttpUploadChk HttpChkInterface
}

func (c *S3Client) initClient() {
	sess := session.Must(session.NewSession())
	cfg := aws.NewConfig()
	cfg.WithEndpoint(c.config.endpoint)
	cfg.WithRegion(c.config.region)
	cfg.WithS3ForcePathStyle(c.config.forcePath)
	creds := credentials.NewStaticCredentials(
		c.config.accessKey,
		c.config.secretKey,
		c.config.token,
	)
	cfg.WithCredentials(creds)
	svc := s3.New(sess, cfg)
	c.client = svc
	c.updriver = s3manager.NewUploaderWithClient(svc)
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

func (c *S3Client) GetConfig() *S3Config {
	return c.config
}

func (c *S3Client) GetUploadDriver() *s3manager.Uploader {
	return c.updriver
}

func (c *S3Client) GetHttpChk() HttpChkInterface {
	if c.HttpUploadChk == nil {
		return defaultHttpChk()
	}
	return c.HttpUploadChk
}

func (c *S3Client) CountObjectInFolder(bucket, folder string) (int, error) {
	if !strings.HasSuffix(folder, "/") {
		folder = folder + "/"
	}
	return c.CountObjectWithPrefix(bucket, folder)
}

func (c *S3Client) CountObjectWithPrefix(bucketName, prefix string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.GetConfig().GetTimeout().GetCountObjTime())
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
	ctx, cancel := context.WithTimeout(context.Background(), c.GetConfig().GetTimeout().GetGetURLTime())
	defer cancel()
	req.SetContext(ctx)
	url, err := req.Presign(expire)
	return url, err
}

func (c *S3Client) GetHeadObject(bucket, objKey string) (*s3.HeadObjectOutput, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.GetConfig().GetTimeout().GetHeadObjTime())
	defer cancel()

	return c.client.HeadObjectWithContext(
		ctx,
		&s3.HeadObjectInput{
			Bucket: &bucket,
			Key:    &objKey,
		})
}

func (c *S3Client) UploadHttpResponse(bucket, objKey string, resp *http.Response) (*UploadResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.config.timeout.upload)
	defer cancel()

	tag := strings.TrimSpace(resp.Header.Get("ETag"))
	t := time.Now().Format("2006-01-02 15:04:05")
	meta := map[string]*string{
		c.GetHttpChk().GetEtag(): &tag,
		"modifiedtime":           &t,
	}

	req := NewUploadRequest(resp.Body, ctx, meta)

	head, err := c.GetHeadObject(bucket, objKey)
	if err != nil {
		return c.UploadSend(bucket, objKey, req)
	}

	exist := c.GetHttpChk().GetChkFn()(resp.Header, head)
	if exist {
		return &UploadResponse{Output: nil, Skip: true}, nil
	}
	return c.UploadSend(bucket, objKey, req)
}

func (c *S3Client) CheckObjectExist(bucket, objKey string) (bool, error) {
	_, err := c.GetHeadObject(bucket, objKey)
	return err == nil, err
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
	if c.GetConfig().GetTimeout().GetUploadTime() != 0 {
		ctx, cancel = context.WithTimeout(context.Background(), c.GetConfig().GetTimeout().GetUploadTime())
		defer cancel()
	} else {
		ctx = context.Background()
	}

	_, err := c.client.PutObjectWithContext(ctx, params)
	return err
}

func (c *S3Client) CreateBucket(bucketName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.GetConfig().GetTimeout().GetCreateBucketTime())
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
	ctx, cancel := context.WithTimeout(context.Background(), c.GetConfig().GetTimeout().GetBucketExistTime())
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

func (c *S3Client) UploadSend(bucket, objKey string, req *UploadRequest) (*UploadResponse, error) {
	input := &s3manager.UploadInput{
		Bucket:   &bucket,
		Key:      &objKey,
		Body:     req.src,
		Metadata: req.meta,
	}

	out, err := upload(c.updriver, req.ctx, input)
	if err != nil {
		return nil, err
	}
	return &UploadResponse{Output: out, Skip: false}, nil
}
