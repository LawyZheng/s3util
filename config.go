package s3util

type S3ConfigInterface interface {
	GetForcePath() bool
	GetTimeout() S3TimeoutInterface
	GetEndpoint() string
	GetRegion() string
	GetAccessKey() string
	GetSecretKey() string
	GetToken() string
}
type S3Config struct {
	ForcePath bool
	Endpoint  string
	AccessKey string
	SecretKey string
	Token     string
	Region    string
	Timeout   S3TimeoutInterface
}

func NewDefaultConfig(endpoint, ak, sk string) *S3Config {
	return &S3Config{
		ForcePath: true,
		Endpoint:  endpoint,
		AccessKey: ak,
		SecretKey: sk,
		Token:     "",
		Region:    "us-east-1",
		Timeout:   NewDefaultTimeout(),
	}
}

func (c *S3Config) GetTimeout() S3TimeoutInterface {
	return c.Timeout
}

func (c *S3Config) GetEndpoint() string {
	return c.Endpoint
}

func (c *S3Config) GetRegion() string {
	return c.Region
}

func (c *S3Config) GetForcePath() bool {
	return c.ForcePath
}

func (c *S3Config) GetAccessKey() string {
	return c.AccessKey
}

func (c *S3Config) GetSecretKey() string {
	return c.SecretKey
}

func (c *S3Config) GetToken() string {
	return c.Token
}
