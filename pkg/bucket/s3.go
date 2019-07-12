package bucket

import (
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/neo9/mongodb-backups/pkg/config"
	"os"
	"path"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

type S3Bucket struct {
	Session *session.Session
	S3 *config.S3
}

func New(s3 *config.S3) *S3Bucket {
	s3Session := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(s3.Region),
	}))

	return &S3Bucket{
		Session: s3Session,
		S3: s3,
	}
}


func (bucket *S3Bucket) Upload(filename string, destFolder string) error {
	uploader := s3manager.NewUploader(bucket.Session)
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket.S3.Name),
		Key:    aws.String(path.Join(destFolder, path.Base(filename))),
		Body:   file,
		ServerSideEncryption: aws.String("AES256"),
	})

	return err
}