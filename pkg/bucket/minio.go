package bucket

import (
	"context"
	"log"
	"os"
	"path"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/neo9/mongodb-backups/pkg/config"
)

type MinioBucket struct {
	Minio *config.Minio
}

func getMinioClient(minioConf *config.Minio) (*minio.Client, error) {
	accessKeyID, accessKeyIDIsDefined := os.LookupEnv("MINIO_ACCESS_KEY_ID")
	secretAccessKey, secretAccessKeyIsDefined := os.LookupEnv("MINIO_SECRET_ACCESS_KEY")

	if !accessKeyIDIsDefined && !secretAccessKeyIsDefined {
		log.Fatalln("You must define MINIO_ACCESS_KEY_ID and MINIO_SECRET_ACCESS_KEY credentials to use Minio connector !")
	}

	minioClient, err := minio.New(minioConf.Host, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: minioConf.SSL,
	})

	if err != nil {
		return nil, err
	}

	return minioClient, nil
}

func NewMinioBucket(minioConf *config.Minio) *MinioBucket {
	minioClient, err := getMinioClient(minioConf)

	minioClient.MakeBucket(context.Background(), minioConf.Name, minio.MakeBucketOptions{Region: minioConf.Region})
	if err != nil {
		exists, errBucketExists := minioClient.BucketExists(context.Background(), minioConf.Name)
		if errBucketExists == nil && exists {
			log.Printf("We already own %s\n", minioConf.Name)
		} else {
			log.Fatalln(err)
		}
	} else {
		log.Printf("Successfully created %s\n", minioConf.Name)
	}

	if err != nil {
		log.Fatalln(err)
	}

	return &MinioBucket{
		Minio: minioConf,
	}
}

func (bucket *MinioBucket) Upload(filename string, destFolder string) error {
	minioClient, err := getMinioClient(bucket.Minio)
	if err != nil {
		return err
	}

	bucketFileName := destFolder + "/" + filename

	_, err = minioClient.FPutObject(context.Background(), bucket.Minio.Name, bucketFileName, filename, minio.PutObjectOptions{})
	if err != nil {
		return err
	}

	log.Println("upload finished")
	return nil
}

func (bucket *MinioBucket) ListFiles(destFolder string) ([]S3File, error) {
	minioClient, err := getMinioClient(bucket.Minio)
	if err != nil {
		return nil, err
	}

	objects := minioClient.ListObjects(context.Background(), bucket.Minio.Name, minio.ListObjectsOptions{
		Recursive: true,
	})

	var files []S3File

	for object := range objects {
		files = append(files, S3File{
			Name: object.Key,
			Etag: object.ETag,
			Size: object.Size,
		})
	}

	return files, nil
}

func (bucket *MinioBucket) DownloadFile(src string) (string, error) {
	minioClient, err := getMinioClient(bucket.Minio)
	if err != nil {
		return "", err
	}

	filename := path.Join("/tmp", path.Base(src))

	err = minioClient.FGetObject(context.Background(), bucket.Minio.Name, src, filename, minio.GetObjectOptions{})
	if err != nil {
		os.Remove(filename)
		return "", err
	}

	return filename, nil
}

func (bucket *MinioBucket) DeleteFile(filename string) error {
	minioClient, err := getMinioClient(bucket.Minio)
	if err != nil {
		return err
	}

	err = minioClient.RemoveObject(context.Background(), bucket.Minio.Name, filename, minio.RemoveObjectOptions{})
	if err != nil {
		log.Printf("unable to delete file %q: %v\n", filename, err)
		return err
	}

	return nil
}
