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
		scheduler.incBackupMetric("error")
		return
	}

	scheduler.addDurationMetric(mongoDBDump.Duration)

	err0 := scheduler.uploadToS3(mongoDBDump.ArchiveFile, scheduler.Plan.Name)
	err1 := scheduler.uploadToS3(mongoDBDump.LogFile, scheduler.Plan.Name)

	if err0 != nil || err1 != nil {
		scheduler.incBackupMetric("error")
	} else {
		scheduler.incBackupMetric("success")
	}
}

func (scheduler *Scheduler) incBackupMetric(status string) {
	scheduler.Metrics.Total.WithLabelValues(scheduler.Plan.Name, status).Inc()
}

func (scheduler *Scheduler) addDurationMetric(duration float64) {
	scheduler.Metrics.Duration.WithLabelValues(scheduler.Plan.Name).Observe(duration)
}

func (scheduler *Scheduler) addSizeMetricFromBackup(filename string) {
	file, err := os.Stat(filename)
	if err != nil {
		log.Errorf("Error computing file size for %s: %v", filename, err)
		return
	}

	scheduler.Metrics.Size.WithLabelValues(scheduler.Plan.Name).Add(float64(file.Size()))
}

func (scheduler *Scheduler) uploadToS3(filename string, destFolder string) error {
	log.Infof("Uploading mongodb file %s", path.Base(filename))

	err := scheduler.Bucket.Upload(filename, destFolder)
	if err != nil {
		log.Errorf("Could not upload to S3: %v", err)
		return err
	}

	err = os.Remove(filename)
	if err != nil {
		log.Errorf("Could not delete file: %v", err)
		return err
	}

	return nil
}

