package bucket

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/neo9/mongodb-backups/pkg/config"
	"os"
	"path"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)


type Bucket interface {
	Upload(filename string, destFolder string) error
	ListFiles(destFolder string) ([]S3File, error)
	DeleteFile(filename string) error
}

type S3Bucket struct {
	Session *session.Session
	S3 *config.S3
}

type S3File struct {
	Etag string
	Name string
	Size int64
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

func (bucket *S3Bucket) ListFiles(destFolder string) ([]S3File, error) {
    svc := s3.New(bucket.Session)

	var files []S3File
    i := 0
	err := svc.ListObjectsPages(&s3.ListObjectsInput{
		Bucket: &bucket.S3.Name,
		Prefix: &destFolder,
	}, func(p *s3.ListObjectsOutput, last bool) (shouldContinue bool) {
		i++

		for _, obj := range p.Contents {
			files = append(files, S3File{
				Name: *obj.Key,
				Etag: strings.ReplaceAll(*obj.ETag, "\"", ""),
				Size: *obj.Size,
			})
		}
		return true
	})
	if err != nil {
		fmt.Println("failed to list objects", err)
		return []S3File{}, err
	}

	return files, nil
}

func (bucket *S3Bucket) DeleteFile(filename string) error {
	svc := s3.New(bucket.Session)
	_, err := svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: &bucket.S3.Name,
		Key: aws.String("//" + filename),
	})

	return err
}
