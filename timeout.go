package s3util

import "time"

type S3Timeout struct {
	Upload       time.Duration
	HeadObject   time.Duration
	CreateBucket time.Duration
	BucketExist  time.Duration
	CountObject  time.Duration
	GetUrl       time.Duration
}

func NewDefaultTimeout() *S3Timeout {
	return &S3Timeout{
		Upload:       time.Minute,
		HeadObject:   3 * time.Second,
		CreateBucket: 3 * time.Second,
		BucketExist:  3 * time.Second,
		CountObject:  10 * time.Second,
		GetUrl:       time.Second,
	}
}

func (t *S3Timeout) GetUploadTime() time.Duration {
	return t.Upload
}

func (t *S3Timeout) GetHeadObjTime() time.Duration {
	return t.HeadObject
}

func (t *S3Timeout) GetCreateBucketTime() time.Duration {
	return t.CreateBucket
}

func (t *S3Timeout) GetBucketExistTime() time.Duration {
	return t.BucketExist
}

func (t *S3Timeout) GetCountObjTime() time.Duration {
	return t.CountObject
}

func (t *S3Timeout) GetGetURLTime() time.Duration {
	return t.GetUrl
}
