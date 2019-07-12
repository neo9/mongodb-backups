package scheduler

import (
	"github.com/neo9/mongodb-backups/pkg/bucket"
	"github.com/neo9/mongodb-backups/pkg/config"
	"github.com/robfig/cron"
)

type Scheduler struct {
	Cron *cron.Cron
	Plan *config.Plan
	S3Bucket *bucket.S3Bucket
}