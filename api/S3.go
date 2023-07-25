package api

import (
	"context"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type contextKey int

const (
	contextKeyS3 contextKey = iota
)

type S3Config struct {
	*s3.S3

	Region    string
	AccessKey string
	SecretKey string
	Buckets   map[string]string
}

func ContextWithS3(ctx context.Context, config *S3Config) context.Context {
	var s3Session *session.Session
	var err error

	s3Session, err = session.NewSession(&aws.Config{
		Region:      aws.String(config.Region),
		Credentials: credentials.NewStaticCredentials(config.AccessKey, config.SecretKey, ""),
	})
	if err != nil {
		panic(fmt.Errorf("failed to connect to s3: %w", err))
	}

	config.S3 = s3.New(s3Session)

	return context.WithValue(ctx, contextKeyS3, config)
}

func S3FromContext(ctx context.Context) *S3Config {
	if rv := ctx.Value(contextKeyS3); rv != nil {
		return rv.(*S3Config)
	}

	return nil
}

func (S3 *S3Config) UploadAttachment(file multipart.File, uuid string, extension string) error {
	putObjectInput := &s3.PutObjectInput{
		Key:  aws.String(uuid + "." + extension),
		Body: file,
	}

	putObjectInput.Bucket = aws.String(S3.Buckets["attachment"])

	_, err := S3.PutObject(putObjectInput)
	if err != nil {
		return ErrUploadError{Err: err, UploadType: "attachment to s3"}
	}

	return nil
}

func (S3 *S3Config) UploadFileToCDN(file multipart.File, uuid string, extension string, fileType string) error {
	putObjectInput := &s3.PutObjectInput{
		Key:         aws.String(uuid + "." + extension),
		Body:        file,
		ContentType: aws.String(fileType),
	}

	putObjectInput.Bucket = aws.String(S3.Buckets["cdn"])

	_, err := S3.PutObject(putObjectInput)
	if err != nil {
		return ErrUploadError{Err: err, UploadType: "file to s3"}
	}

	return nil
}

func (S3 *S3Config) GetAttachmentURL(uuid string, extension string, filename string) (string, error) {

	getObjectInput := &s3.GetObjectInput{
		Key:                        aws.String(uuid + "." + extension),
		ResponseContentDisposition: aws.String("attachment; filename=" + filename),
	}

	getObjectInput.Bucket = aws.String(S3.Buckets["attachment"])

	req, _ := S3.GetObjectRequest(getObjectInput)
	url, err := req.Presign(time.Hour * 24)
	if err != nil {
		return "", err
	}

	return url, nil
}

func (S3 *S3Config) DeleteAttachmentFromS3(uuid string, extension string) error {

	deleteObjectInput := &s3.DeleteObjectInput{
		Key: aws.String(uuid + extension),
	}

	deleteObjectInput.Bucket = aws.String(S3.Buckets["attachment"])

	_, err := S3.DeleteObject(deleteObjectInput)
	if err != nil {
		return err
	}

	return nil
}
