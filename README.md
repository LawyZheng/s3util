## s3util

s3util is an easy way to interact with S3 or S3-compatible storage.


### Installation
```shell
go get github.com/lawyzheng/s3util
```

### Quick Start

- #### Create Client
```go
client := s3util.NewS3Client("my_endpoint", "my_accesskey", "my_scretekey")
```

- #### Create Bucket
```go
err := client.CreateBucket("my-bucket")
if err != nil {
    // handler error
    ...
}
```

- #### Get Object Download URL
```go
urlString, err := client.GetObjectDownloadURL("my-bucket", "my_key", time.Minute)
if err != nil {
    // handler error
    ...
}
```

- #### Upload Object From HTTP Response
```go
resp, err := http.Get("my_resource_url")
if err != nil {
    // handler error
    ...
}
defer resp.Body.Close()

if resp.StatusCode != http.StatusOK {
    // handler status code which is not 200
    ...
}

// skip will be true, if the object have already existed.
skip, err := client.UploadHttpResponse("my-bucket", "my_key", resp)
if err != nil {
    // handle error
    ...
}

```


