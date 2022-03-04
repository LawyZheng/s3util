package s3util

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUploadReader(t *testing.T) {
	assert := assert.New(t)

	data := []byte("12345678")
	reader := bytes.NewReader(data)
	bucket := "test-upload"
	key := "test.txt"
	c := NewS3Client("http://192.168.0.52:9000", "kayisoftadmin", "kayisoftadmin")

	out, err := c.UploadReader(bucket, key, reader)
	assert.Equal(nil, err, "response error should be nil")
	fmt.Printf("out: %v\n", out)

	key = "test.dcm"
	f, _ := os.Open(key)
	out, err = c.UploadReader(bucket, key, f)
	assert.Equal(nil, err, "response error should be nil")
	fmt.Printf("out: %v\n", out)
}
