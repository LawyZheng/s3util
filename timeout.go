package s3util

import "time"

type S3Timeout struct {
	upload       time.Duration
	headObject   time.Duration
	createBucket time.Duration
	bucketExist  time.Duration
	countObject  time.Duration
	getUrl       time.Duration
}

func NewDefaultTimeout() *S3Timeout {
	return &S3Timeout{
		upload:       time.Minute,
		headObject:   3 * time.Second,
		createBucket: 3 * time.Second,
		bucketExist:  3 * time.Second,
		countObject:  10 * time.Second,
		getUrl:       time.Second,
	}
}

func (t *S3Timeout) GetUploadTime() time.Duration {
	return t.upload
}

func (t *S3Timeout) SetUploadTime(timeout time.Duration) {
	t.upload = timeout
}

func (t *S3Timeout) GetHeadObjTime() time.Duration {
	return t.headObject
}

func (t *S3Timeout) SetHeadObjTime(timeout time.Duration) {
	t.headObject = timeout
}

func (t *S3Timeout) GetCreateBucketTime() time.Duration {
	return t.createBucket
}

func (t *S3Timeout) SetCreateBucketTime(timeout time.Duration) {
	t.createBucket = timeout
}

func (t *S3Timeout) GetBucketExistTime() time.Duration {
	return t.bucketExist
}

func (t *S3Timeout) SetBucketExistTime(timeout time.Duration) {
	t.bucketExist = timeout
}

func (t *S3Timeout) GetCountObjTime() time.Duration {
	return t.countObject
}

func (t *S3Timeout) SetCountObjTime(timeout time.Duration) {
	t.countObject = timeout
}

func (t *S3Timeout) GetGetURLTime() time.Duration {
	return t.getUrl
}

func (t *S3Timeout) SetGetURLTime(timeout time.Duration) {
	t.getUrl = timeout
}
