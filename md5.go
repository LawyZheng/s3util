package s3util

import (
	"crypto/md5"
	"encoding/hex"
	"io"
)

func GetMd5(src io.ReadSeeker) string {
	md5h := md5.New()
	n, _ := io.Copy(md5h, src)
	src.Seek(-n, io.SeekCurrent)
	return hex.EncodeToString(md5h.Sum(nil))
}
