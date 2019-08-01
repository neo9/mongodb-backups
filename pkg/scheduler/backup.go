package scheduler

import (
	"github.com/neo9/mongodb-backups/pkg/mongodb"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
}

func (scheduler *Scheduler) runBackup() {
	log.Infof("Running mongodb %s", scheduler.Plan.Name)

	mongoDBDump, err := mongodb.CreateDump(scheduler.Plan)
	if err != nil {
		log.Errorf("Error creating dump for %s", scheduler.Plan.Name)
		return
	}

	scheduler.uploadToS3(mongoDBDump.ArchiveFile, scheduler.Plan.Name)
	scheduler.uploadToS3(mongoDBDump.LogFile, scheduler.Plan.Name)
}

func (scheduler *Scheduler) uploadToS3(filename string, destFolder string) {
	log.Infof("Uploading mongodb file %s", path.Base(filename))

	err := scheduler.Bucket.Upload(filename, destFolder)
	if err != nil {
		log.Errorf("Could not upload to S3: %v", err)
	}

	err = os.Remove(filename)
	if err != nil {
		log.Errorf("Could not delete file: %v", err)
	}
}

