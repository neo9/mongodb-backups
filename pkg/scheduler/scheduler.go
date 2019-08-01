package scheduler

import (
	"fmt"
	"github.com/neo9/mongodb-backups/pkg/bucket"
	"github.com/neo9/mongodb-backups/pkg/config"
	"github.com/neo9/mongodb-backups/pkg/mongodb"
	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
)

type Scheduler struct {
	Cron *cron.Cron
	Plan *config.Plan
	Bucket bucket.Bucket
}

func init() {
	log.SetFormatter(&log.JSONFormatter{})
}

func New(plan *config.Plan) *Scheduler {
	S3Bucket := bucket.New(&plan.Bucket.S3)

	return &Scheduler{
		Plan: plan,
		Cron: cron.New(),
		Bucket: S3Bucket,
	}
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

func (scheduler *Scheduler) Run() {
	err := scheduler.Cron.AddFunc(fmt.Sprintf("0 %s", scheduler.Plan.Schedule), func() {
		scheduler.runBackup()
	})

	if err != nil {
		log.Errorf("Could not schedule mongodb %s, error: %v", scheduler.Plan.Name, err)
	}

	log.Infof("Name: %s, Schedule: %s", scheduler.Plan.Name, scheduler.Plan.Schedule)

	scheduler.Cron.Start()
	scheduler.displaySchedule()

	// TODO: remove
	scheduler.runBackup()
}

func (scheduler *Scheduler) displaySchedule() {
	entries := scheduler.Cron.Entries()

	entry := entries[0]
	log.Infof("Backup %s will run at %v", scheduler.Plan.Name, entry.Next)
}