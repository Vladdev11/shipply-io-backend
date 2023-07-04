package api

import (
	"errors"
	"mime/multipart"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/shipply-io/shipply-io-backend/util"
)

var S3 *s3.S3

func InitAWSS3() {

	var s3Session *session.Session
	var err error

	s3Session, err = session.NewSession(&aws.Config{
		Region:      aws.String(util.ConfigS3Region),
		Credentials: credentials.NewStaticCredentials(util.ConfigAWSAccessKeyID, util.ConfigAWSSecretAccessKey, ""),
	})
	if err != nil {
		panic(errors.New("failed to connect to s3"))
	}

	S3 = s3.New(s3Session)
}

func UploadAttachmentToS3(file multipart.File, uuid string, extension string) error {
	putObjectInput := &s3.PutObjectInput{
		Key:  aws.String(uuid + "." + extension),
		Body: file,
	}

	putObjectInput.Bucket = aws.String(util.ConfigS3AttachmentBucket)

	_, err := S3.PutObject(putObjectInput)
	if err != nil {
		return err
	}

	return nil
}

func UploadFileToCDN(file multipart.File, uuid string, extension string, fileType string) error {
	putObjectInput := &s3.PutObjectInput{
		Key:         aws.String(uuid + "." + extension),
		Body:        file,
		ContentType: aws.String(fileType),
	}

	putObjectInput.Bucket = aws.String(util.ConfigS3CDNBucket)

	_, err := S3.PutObject(putObjectInput)
	if err != nil {
		return err
	}

	return nil
}

func GetAttachmentURL(uuid string, extension string, filename string) (string, error) {

	getObjectInput := &s3.GetObjectInput{
		Key:                        aws.String(uuid + "." + extension),
		ResponseContentDisposition: aws.String("attachment; filename=" + filename),
	}

	getObjectInput.Bucket = aws.String(util.ConfigS3AttachmentBucket)

	req, _ := S3.GetObjectRequest(getObjectInput)
	url, err := req.Presign(time.Hour * 24)
	if err != nil {
		return "", err
	}

	return url, nil
}

func DeleteAttachmentFromS3(uuid string, extension string) error {

	deleteObjectInput := &s3.DeleteObjectInput{
		Key: aws.String(uuid + extension),
	}

	deleteObjectInput.Bucket = aws.String(util.ConfigS3AttachmentBucket)

	_, err := S3.DeleteObject(deleteObjectInput)
	if err != nil {
		return err
	}

	return nil
}
