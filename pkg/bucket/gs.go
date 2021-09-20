package bucket

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/neo9/mongodb-backups/pkg/config"
	"google.golang.org/api/iterator"

	"cloud.google.com/go/storage"

	log "github.com/sirupsen/logrus"
)

type GSBucket struct {
	GS *config.GS
}

func NewGSBucket(gs *config.GS) *GSBucket {
	return &GSBucket{
		GS: gs,
	}
}

func (bucket *GSBucket) Upload(filename string, destFolder string) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()
	bucketClient := client.Bucket(bucket.GS.Name)

	wc := bucketClient.Object(path.Join(destFolder, path.Base(filename))).NewWriter(ctx)
	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	if _, err := io.Copy(wc, file); err != nil {
		return err
	}
	if err := wc.Close(); err != nil {
		return err
	}

	log.Infof("upload finished")
	return nil
}

func (bucket *GSBucket) ListFiles(destFolder string) ([]S3File, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Infof("failed to create client: %v", err)
		return []S3File{}, err
	}
	defer client.Close()
	bucketClient := client.Bucket(bucket.GS.Name)

	var files []S3File
	query := &storage.Query{Prefix: destFolder}
	it := bucketClient.Objects(ctx, query)
	for {
		obj, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			_ = fmt.Errorf("listBucket: unable to list bucket %q: %v", bucket.GS.Name, err)
			return []S3File{}, err
		}
		if destFolder == obj.Name || (destFolder+"/") == obj.Name {
			continue
		}

		files = append(files, S3File{
			Name: obj.Name,
			Etag: strings.ReplaceAll(obj.Etag, "\"", ""),
			Size: obj.Size,
		})
	}

	return files, nil
}

func (bucket *GSBucket) DownloadFile(src string) (string, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Infof("failed to create client: %v", err)
		return "", err
	}
	defer client.Close()
	bucketClient := client.Bucket(bucket.GS.Name)

	filename := path.Join("/tmp", path.Base(src))
	file, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	rc, err := bucketClient.Object(src).NewReader(ctx)
	if err != nil {
		_ = os.Remove(filename)
		return "", err
	}

	_, err = io.Copy(file, rc)
	if err != nil {
		return "", fmt.Errorf("ioutil.read: %v", err)
	}

	return filename, nil
}

func (bucket *GSBucket) DeleteFile(filename string) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Infof("failed to create client: %v", err)
		return err
	}
	defer client.Close()
	bucketClient := client.Bucket(bucket.GS.Name)

	if err := bucketClient.Object(filename).Delete(ctx); err != nil {
		log.Errorf("unable to delete file %q: %v", filename, err)
		return err
	}

	return nil
}
