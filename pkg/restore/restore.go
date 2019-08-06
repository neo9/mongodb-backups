package restore

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/neo9/mongodb-backups/pkg/bucket"
	"github.com/neo9/mongodb-backups/pkg/scheduler"
	"github.com/neo9/mongodb-backups/pkg/utils"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

func DisplayBackups(scheduler *scheduler.Scheduler) error {
	files, err := getFiles(scheduler)
	if err != nil {
		return err
	}

	if len(files) == 0 {
		log.Infof("Bucket is empty")
		return nil
	}

	for i := 0; i < len(files); i++ {
		file := files[i]
		if strings.Contains(file.Name, ".log") {
			continue
		}

		timestamp, err := utils.GetBucketFileTimestamp(file.Name)
		if err != nil {
			log.Errorf("Could not parse file: %v", err)
			continue
		}

		log.Infof("%s | Backup %s, size: %s",
			time.Unix(timestamp, 10),
			file.Etag,
			utils.GetHumanBytes(file.Size))
	}

	return nil
}

func Restore(scheduler *scheduler.Scheduler, restoreID string) error {
	files, err := getFiles(scheduler)
	if err != nil {
		return err
	}

	var file bucket.S3File
	for i := 0; i < len(files); i++ {
		if *aws.String(restoreID) == files[i].Etag {
			file = files[i]
			break
		}
	}

	if file.Name == "" {
		log.Errorf("Could not find tag %s in folder %s", restoreID, scheduler.Plan.Name)
		return errors.New("TAG_NOT_FOUND")
	}

	log.Infof("Restoring backup %s from snapshot %s", restoreID, file.Name)

	return nil
}

func downloadBackup(s3Path string) error {
	// TODO: download S3 file

	// Apply default host

	// DELETE tmp
	return nil
}

func getFiles(scheduler *scheduler.Scheduler) ([]bucket.S3File, error) {
	files, err := scheduler.Bucket.ListFiles(scheduler.Plan.Name)
	if err != nil {
		log.Errorf("Could not list files for %s: %v", scheduler.Plan.Name, err)
		return []bucket.S3File{}, err
	}

	return files, nil
}

