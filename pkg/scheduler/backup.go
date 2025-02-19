package scheduler

import (
	"os"
	"path"
	"time"

	"github.com/neo9/mongodb-backups/pkg/log"
	"github.com/neo9/mongodb-backups/pkg/mongodb"
)

func (scheduler *Scheduler) runBackup() {
	log.Info("Running mongodb %s", scheduler.Plan.Name)

	mongoDBDump, err := scheduler.runDump()
	if err != nil {
		log.Error("Error creating dump for %s", scheduler.Plan.Name)
		scheduler.incTotalBackupMetricWithStatus("error")
		return
	}

	scheduler.addMongoDumpLatencyMetric(mongoDBDump.Duration)
	scheduler.addBackupSizMetric(mongoDBDump.ArchiveFile)

	err0 := scheduler.uploadToS3(mongoDBDump.ArchiveFile, scheduler.Plan.Name)
	err1 := scheduler.uploadToS3(mongoDBDump.LogFile, scheduler.Plan.Name)

	if err0 != nil || err1 != nil {
		scheduler.incTotalBackupMetricWithStatus("error")
	} else {
		scheduler.incTotalBackupMetricWithStatus("success")
		timestamp := float64(time.Now().Unix())
		scheduler.Metrics.LastSuccessfulSnapshot.WithLabelValues(scheduler.Plan.Name).Set(timestamp)
	}
}

func (scheduler *Scheduler) runDump() (mongodb.MongoDBDump, error) {
	var err error
	var mongoDBDump mongodb.MongoDBDump
	maxRetries := scheduler.Plan.CreateDump.MaxRetries
	retryDelay := scheduler.Plan.CreateDump.RetryDelay * time.Second

	for i := 0; i < maxRetries; i++ {
		mongoDBDump, err := mongodb.CreateDump(scheduler.Plan)
		if err != nil {
			scheduler.incMongodbBackupTries("error")
			log.Error("Error creating mongodump (retry %d/%d): %v", i+1, maxRetries, err)
			time.Sleep(retryDelay)
		} else {
			scheduler.incMongodbBackupTries("success")
			return mongoDBDump, nil
		}
	}
	return mongoDBDump, err

}

func (scheduler *Scheduler) incTotalBackupMetricWithStatus(status string) {
	scheduler.Metrics.Total.WithLabelValues(scheduler.Plan.Name, status).Inc()
}

func (scheduler *Scheduler) addMongoDumpLatencyMetric(duration float64) {
	scheduler.Metrics.SnapshotLatency.WithLabelValues(scheduler.Plan.Name).Observe(duration)
}

func (scheduler *Scheduler) addBackupSizMetric(filename string) {
	file, err := os.Stat(filename)
	if err != nil {
		log.Error("Error computing file size for %s: %v", filename, err)
		return
	}

	scheduler.Metrics.BackupSize.WithLabelValues(scheduler.Plan.Name).Add(float64(file.Size()))
}

func (scheduler *Scheduler) uploadToS3(filename string, destFolder string) error {
	log.Info("Uploading mongodb file %s", path.Base(filename))

	err := scheduler.Bucket.Upload(filename, destFolder)
	if err != nil {
		log.Error("Could not upload to S3: %v", err)
		return err
	}

	err = os.Remove(filename)
	if err != nil {
		log.Error("Could not delete file: %v", err)
		return err
	}

	return nil
}

func (scheduler *Scheduler) incMongodbBackupTries(status string) {
	scheduler.Metrics.BackupTries.WithLabelValues(scheduler.Plan.Name, status).Inc()
}
