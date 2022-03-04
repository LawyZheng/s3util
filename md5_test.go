package s3util

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMd5(t *testing.T) {
	assert := assert.New(t)

	s := []byte("12345678")
	reader := bytes.NewReader(s)

	target := make([]byte, 2)
	reader.Read(target)
	assert.Equal("12", string(target), "Read Two Bytes")
	assert.Equal("5bd2026f128662763c532f2f4b6f2476", GetMd5(reader), "MD5 345678")

	b, _ := io.ReadAll(reader)
	assert.Equal("345678", string(b), "Read After MD5")
}

func TestFileMd5(t *testing.T) {
	f, _ := os.Open("test.dcm")
	fmt.Printf("getMd5(f): %v\n", GetMd5(f))
}
