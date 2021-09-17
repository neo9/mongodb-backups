package scheduler

import (
	"fmt"
	"time"

	"github.com/neo9/mongodb-backups/pkg/utils"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
}

func (scheduler *Scheduler) deleteOldBackups() {
	files, err := scheduler.Bucket.ListFiles(scheduler.Plan.Name)
	if err != nil {
		scheduler.incRetentionMetricError(fmt.Sprintf("Could not list files for plan %s", scheduler.Plan.Name))
		return
	}

	retentionDuration, err := utils.GetDurationFromTimeString(scheduler.Plan.Retention)
	if err != nil {
		scheduler.incRetentionMetricError(fmt.Sprintf("Could not execute retention: %v", err))
		return
	}

	var removeFiles []string

	for i := 0; i < len(files); i++ {
		file := files[i]
		log.Debugf("File: ", file.Name)
		timestamp, err := utils.GetBucketFileTimestamp(file.Name)
		if err != nil {
			scheduler.incRetentionMetricError(fmt.Sprintf("Could not apply retention: %v", err))
		}

		ageInSeconds := time.Now().Unix() - timestamp
		diffInSeconds := ageInSeconds - int64(retentionDuration.Seconds())

		if diffInSeconds > 0 {
			log.Debugf("File is %s old and schedule for removal", file.Name)
			removeFiles = append(removeFiles, file.Name)
		} else {
			log.Debugf("File is %s old and will be removed in %s", file.Name,
				time.Duration(diffInSeconds*-1)*time.Second)
		}
	}

	log.Infof("Retention: %d file(s) to remove", len(removeFiles))

	status := "success"
	for i := 0; i < len(removeFiles); i++ {
		err := scheduler.Bucket.DeleteFile(removeFiles[i])
		if err != nil {
			scheduler.incRetentionMetricError(fmt.Sprintf("Could not remove file %s", removeFiles[i]))
			status = "error"
		}
	}

	snapshotCount := float64(len(files)/2 - len(removeFiles))
	scheduler.Metrics.BucketCount.WithLabelValues(scheduler.Plan.Name).Set(snapshotCount)
	scheduler.Metrics.RetentionTotal.WithLabelValues(scheduler.Plan.Name, status).Inc()
}

func (scheduler *Scheduler) incRetentionMetricError(error string) {
	log.Error(error)
	scheduler.Metrics.RetentionTotal.WithLabelValues(scheduler.Plan.Name, "error").Inc()
}
