package storage

import (
    "bytes"
    "context"
    "fmt"
    "io"
    
    "github.com/Brownie44l1/blog/config"
    "github.com/aws/aws-sdk-go-v2/aws"
    awsConfig "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/credentials"
    "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type S3Storage struct {
    client *s3.Client
    bucket string
}

func New(cfg *config.Config) (*s3.Client, error) {
    awsCfg, err := awsConfig.LoadDefaultConfig(
        context.TODO(),
        awsConfig.WithRegion(cfg.S3Region),
        awsConfig.WithCredentialsProvider(
            credentials.NewStaticCredentialsProvider(
                cfg.AWSAccessKey,
                cfg.AWSSecretKey,
                "",
            ),
        ),
    )
    if err != nil {
        return nil, err
    }

    return s3.NewFromConfig(awsCfg), nil
}

func UploadFile(s3c *s3.Client, bucket, key string, data []byte) error {
    _, err := s3c.PutObject(context.TODO(), &s3.PutObjectInput{
        Bucket: aws.String(bucket),
        Key:    aws.String(key),
        Body:   bytes.NewReader(data),
    })
    return err
}

func (s *S3Storage) UploadImage(blogID int64, file io.Reader, filename string) (string, error) {
    key := fmt.Sprintf("blogs/%d/%s", blogID, filename)
    
    // Read file into buffer
    buf := new(bytes.Buffer)
    if _, err := io.Copy(buf, file); err != nil {
        return "", err
    }
    
    _, err := s.client.PutObject(context.TODO(), &s3.PutObjectInput{
        Bucket: aws.String(s.bucket),
        Key:    aws.String(key),
        Body:   bytes.NewReader(buf.Bytes()),
        ACL:    types.ObjectCannedACLPublicRead,
    })
    
    if err != nil {
        return "", err
    }
    
    url := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", s.bucket, key)
    return url, nil
}