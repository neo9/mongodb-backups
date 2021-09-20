package restore

import (
	"errors"
	"os"
	"path"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/neo9/mongodb-backups/pkg/bucket"
	"github.com/neo9/mongodb-backups/pkg/mongodb"
	"github.com/neo9/mongodb-backups/pkg/scheduler"
	"github.com/neo9/mongodb-backups/pkg/utils"
	log "github.com/sirupsen/logrus"
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
		timestamp, err := utils.GetBucketFileTimestamp(path.Base(file.Name))
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

func RestoreLast(scheduler *scheduler.Scheduler, args string) error {
	files, err := getFiles(scheduler)
	if err != nil {
		return err
	}

	if len(files) == 0 {
		log.Error("No backup found")
		return errors.New("NO_BACKUP")
	}

	file := files[len(files)-1]
	log.Infof("Restoring backup %s from snapshot %s", file.Etag, file.Name)
	return restoreBackup(scheduler, file.Name, args)
}

func Restore(scheduler *scheduler.Scheduler, restoreID string, args string) error {
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
	return restoreBackup(scheduler, file.Name, args)
}

func restoreBackup(scheduler *scheduler.Scheduler, src string, args string) error {
	log.Info("Download snapshot")
	file, err := scheduler.Bucket.DownloadFile(src)
	if err != nil {
		log.Errorf("Could not download snapshot: %v", err)
		_ = os.Remove(file)
		return err
	}

	err = mongodb.RestoreDump(file, args, scheduler.Plan)
	_ = os.Remove(file)
	return err
}

func getFiles(scheduler *scheduler.Scheduler) ([]bucket.S3File, error) {
	files, err := scheduler.Bucket.ListFiles(scheduler.Plan.Name)
	if err != nil {
		log.Errorf("Could not list files for %s: %v", scheduler.Plan.Name, err)
		return []bucket.S3File{}, err
	}

	var dumpFiles []bucket.S3File
	for i := 0; i < len(files); i++ {
		if strings.Contains(files[i].Name, ".gz") {
			dumpFiles = append(dumpFiles, files[i])
		}
	}

	return dumpFiles, nil
}
