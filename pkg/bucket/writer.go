package bucket

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/neo9/mongodb-backups/pkg/utils"
	"io"
	"sync/atomic"
)

type progressWriter struct {
	written int64
	writer  io.WriterAt
	size    int64
	humanSize string
}

func (pw *progressWriter) WriteAt(p []byte, off int64) (int, error) {
	atomic.AddInt64(&pw.written, int64(len(p)))

	percentageDownloaded := float32(pw.written*100) / float32(pw.size)

	fmt.Printf("File size: %s downloaded: %s percentage: %.2f%%\r", pw.humanSize, utils.GetHumanBytes(pw.written), percentageDownloaded)

	return pw.writer.WriteAt(p, off)
}

func getFileSize(svc *s3.S3, bucket string, prefix string) (filesize int64, error error) {
	params := &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(prefix),
	}

	resp, err := svc.HeadObject(params)
	if err != nil {
		return 0, err
	}

	return *resp.ContentLength, nil
}

