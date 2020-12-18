package bucket

import (
	"github.com/neo9/mongodb-backups/pkg/config"
	log "github.com/sirupsen/logrus"
)

type Bucket interface {
	Upload(filename string, destFolder string) error
	ListFiles(destFolder string) ([]S3File, error)
	DownloadFile(src string) (string, error)
	DeleteFile(filename string) error
}

func New(bucket *config.Bucket) Bucket {
	if bucket.S3.Name != "" {
		log.Infof("using s3 storage (%s)", bucket.S3.Name)
		return NewS3Bucket(&bucket.S3)
	}
	if bucket.GS.Name != "" {
		log.Infof("using gs storage (%s)", bucket.GS.Name)
		return NewGSBucket(&bucket.GS)
	}

	panic("No implementations declared for bucket configuration")
}
