package scheduler

import (
	"fmt"

	"github.com/neo9/mongodb-backups/pkg/bucket"
	"github.com/neo9/mongodb-backups/pkg/config"
	"github.com/neo9/mongodb-backups/pkg/log"
	"github.com/neo9/mongodb-backups/pkg/metrics"
	"github.com/robfig/cron"
)

type Scheduler struct {
	Cron    *cron.Cron
	Plan    *config.Plan
	Bucket  bucket.Bucket
	Metrics *metrics.BackupMetrics
}

func New(plan *config.Plan) *Scheduler {
	Bucket := bucket.New(&plan.Bucket)

	return &Scheduler{
		Plan:    plan,
		Cron:    cron.New(),
		Bucket:  Bucket,
		Metrics: metrics.New("mongodb_backups", "scheduler"),
	}
}

func (scheduler *Scheduler) Run() {
	err := scheduler.Cron.AddFunc(fmt.Sprintf("0 %s", scheduler.Plan.Schedule), func() {
		scheduler.runBackup()
	})
	if err != nil {
		log.Error("Could not schedule mongodb %s, error: %v", scheduler.Plan.Name, err)
	}

	log.Info("Name: %s, Schedule: %s", scheduler.Plan.Name, scheduler.Plan.Schedule)
	log.Info("MaxRetries: %s, RetryDelay: %s", scheduler.Plan.CreateDump.MaxRetries, scheduler.Plan.CreateDump.RetryDelay)

	err = scheduler.Cron.AddFunc("0 0 * * * *", func() {
		scheduler.DeleteOldBackups()
	})
	if err != nil {
		log.Error("Could not schedule retention, error: %v", err)
	}

	scheduler.Cron.Start()
	scheduler.displaySchedule()
}

func (scheduler *Scheduler) displaySchedule() {
	entries := scheduler.Cron.Entries()

	entry := entries[0]
	log.Info("Backup %s will run at %v", scheduler.Plan.Name, entry.Next)
}
