package s3util

type S3Config struct {
	forcePath bool
	endpoint  string
	accessKey string
	secretKey string
	token     string
	region    string
	timeout   *S3Timeout
}

func NewDefaultConfig(endpoint, ak, sk string) *S3Config {
	return &S3Config{
		forcePath: true,
		endpoint:  endpoint,
		accessKey: ak,
		secretKey: sk,
		region:    "us-east-1",
		timeout:   NewDefaultTimeout(),
	}
}

func (c *S3Config) SetTimeout(timeout *S3Timeout) {
	c.timeout = timeout
}

func (c *S3Config) GetTimeout() *S3Timeout {
	return c.timeout
}

func (c *S3Config) SetRegion(region string) {
	c.region = region
}

func (c *S3Config) GetRegion() string {
	return c.region
}
