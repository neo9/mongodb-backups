package scheduler

import (
	"fmt"
	mongodb "github.com/neo9/mongodb-backups/pkg/backup"
	"github.com/neo9/mongodb-backups/pkg/bucket"
	"github.com/neo9/mongodb-backups/pkg/config"
	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"
	"os"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
}

func New(plan *config.Plan) *Scheduler {
	return &Scheduler{
		Plan: plan,
		Cron: cron.New(),
		S3Bucket: bucket.New(&plan.S3),
	}
}

func (scheduler *Scheduler) runBackup(backup config.Backup) {
	log.Infof("Running backup %s", backup.Name)
	filenames, err := mongodb.MongoDBDump(&backup.MongoDB, backup.Name)
	if err != nil {
		log.Errorf("Error creating dump for %s", backup.Name)
	}

	for i := 0; i < len(filenames); i++ {
		log.Infof("Uploading file %s to bucket %s", filenames[i], scheduler.S3Bucket.S3.Name)
		err := scheduler.S3Bucket.Upload(filenames[i], backup.MongoDB.Host)
		if err != nil {
			log.Errorf("Could not upload to S3: %v", err)
		}
		err = os.Remove(filenames[i])
		if err != nil {
			log.Errorf("Could not delete log file: %v", err)
		}
	}
}

func (scheduler *Scheduler) LaunchScheduler() {
	for i := 0; i < len(scheduler.Plan.Backups); i++ {
		backup := scheduler.Plan.Backups[i]
		err := scheduler.Cron.AddFunc(fmt.Sprintf("0 %s", backup.Schedule), func() {
			scheduler.runBackup(backup)
		})
		if err != nil {
			log.Errorf("Could schedule backup %s, error: %v", backup.Name, err)
		}

		log.Infof("Name: %s, Schedule: %s", backup.Name, backup.Schedule)
	}

	scheduler.Cron.Start()
	scheduler.displaySchedules()
	scheduler.runBackup(scheduler.Plan.Backups[0])
}

func (scheduler *Scheduler) displaySchedules() {
	entries := scheduler.Cron.Entries()
	if len(scheduler.Plan.Backups) != len(entries) {
		panic("Backup and cron entries are not the same length")
	}

	for i := 0; i < len(entries); i++ {
		entry := entries[i]
		backup := scheduler.Plan.Backups[i]
		log.Infof("Backup %s will run at %v", backup.Name, entry.Next)
	}
}