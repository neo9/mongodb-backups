package scheduler

import (
	"fmt"
	"github.com/neo9/mongodb-backups/pkg/bucket"
	"github.com/neo9/mongodb-backups/pkg/config"
	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"
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


func (scheduler *Scheduler) Run() {
	err := scheduler.Cron.AddFunc(fmt.Sprintf("0 %s", scheduler.Plan.Schedule), func() {
		scheduler.runBackup()
	})
	if err != nil {
		log.Errorf("Could not schedule mongodb %s, error: %v", scheduler.Plan.Name, err)
	}

	log.Infof("Name: %s, Schedule: %s", scheduler.Plan.Name, scheduler.Plan.Schedule)

	err = scheduler.Cron.AddFunc("0 0 * * * *", func() {
		scheduler.deleteOldBackups()
	})
	if err != nil {
		log.Errorf("Could not schedule retention, error: %v", err)
	}

	scheduler.Cron.Start()
	scheduler.displaySchedule()

	// TODO: remove
	scheduler.deleteOldBackups()
}

func (scheduler *Scheduler) displaySchedule() {
	entries := scheduler.Cron.Entries()

	entry := entries[0]
	log.Infof("Backup %s will run at %v", scheduler.Plan.Name, entry.Next)
}