package s3util

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func upload(driver *s3manager.Uploader, ctx context.Context, input *s3manager.UploadInput) (*s3manager.UploadOutput, error) {
	output, err := driver.UploadWithContext(ctx, input)
	if err != nil {
		return output, fmt.Errorf("UploadFailed: caused by: %s", err)
	}
	return output, nil
}
