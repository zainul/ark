package xs3

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// S3Upload type of s3
type S3Upload struct {
	Bucket      string // S3 Bucket
	BaseFolder  string // Base folder for tidying up files by folder
	Region      string // S3 region
	s3Client    *s3.S3
	Credentials *credentials.Credentials
	MaxRetry    int
}

// NewS3Upload is Constructor for S3Upload struct
func NewS3Upload(bucket string, baseFolder string, accessKey string, secretKey string, region string, maxRetry int) *S3Upload {
	// Initialize credentials. Currently we use static credentials, and token is empty string
	// since we are not requiring it.
	cred := credentials.NewStaticCredentials(accessKey, secretKey, "")

	// Session must be created once and can be reused. According to AWS docs,
	// this session is thread safe.
	// This session configured with 2 times retry if timeout happened.
	awsSession := session.Must(session.NewSession(&aws.Config{
		Credentials: cred,
		Region:      aws.String(region),
		MaxRetries:  aws.Int(maxRetry),
	}))

	return &S3Upload{
		Bucket:      bucket,
		BaseFolder:  baseFolder,
		Region:      region,
		s3Client:    s3.New(awsSession),
		Credentials: cred,
	}
}

// UploadFile is func to upload static file to s3
func (su *S3Upload) UploadFile(key string, body []byte) error {
	timeout := 10 * time.Second

	ctx := context.Background()
	var cancelFn func()
	if timeout > 0 {
		ctx, cancelFn = context.WithTimeout(ctx, timeout)
	}

	// Ensure go routines canceled if timeout
	defer cancelFn()

	_, err := su.s3Client.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket: aws.String(su.Bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(body),
	})

	if err != nil {
		if val, ok := err.(awserr.Error); ok && val.Code() == request.CanceledErrorCode {
			fmt.Fprintf(os.Stderr, "upload canceled due to timeout,  %v\n", err)
		} else {
			fmt.Fprintf(os.Stderr, "failed to upload object, %v\n", err)
		}

		return err
	}

	return nil
}
