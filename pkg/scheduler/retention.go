package scheduler

import (
	"fmt"
	"github.com/neo9/mongodb-backups/pkg/utils"
	log "github.com/sirupsen/logrus"
	"regexp"
	"strconv"
	"time"
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
    	log.Debugf("File: ", files[i])

    	reg := regexp.MustCompile(`mongodb-snapshot-(?P<Time>\d+)\.(gz|log)`)
    	match := reg.FindStringSubmatch(files[i])
    	if len(match) != 3 {
			scheduler.incRetentionMetricError(fmt.Sprintf(
				"File does not match pattern in folder %s: %s", scheduler.Plan.Name, files[i]))
		}

        timestamp, err := strconv.ParseInt(match[1], 10, 64)
        if err != nil {
			scheduler.incRetentionMetricError(fmt.Sprintf(
				"File has invalid timestamp in folder %s: %s", scheduler.Plan.Name, files[i]))
		}

        ageInSeconds := time.Now().Unix() - timestamp
        diffInSeconds := ageInSeconds - int64(retentionDuration.Seconds())

        if diffInSeconds > 0 {
			log.Debugf("File is %s old and schedule for removal", files[i])
        	removeFiles = append(removeFiles, files[i])
		} else {
			log.Debugf("File is %s old and will be removed in %s", files[i],
				time.Duration(diffInSeconds * -1) * time.Second)
		}
	}

    log.Infof("Retention: %d file(s) to remove", len(removeFiles))

    status := "success"
    for i := 0; i < len(removeFiles); i++ {
    	err := scheduler.Bucket.DeleteFile(removeFiles[i])
    	if err != nil {
    		log.Errorf("Could not remove file %s", removeFiles[i])
            status = "error"
		}
	}


	scheduler.Metrics.RetentionTotal.WithLabelValues(scheduler.Plan.Name, status).Inc()
}

func (scheduler *Scheduler) incRetentionMetricError(error string) {
    log.Error(error)
    scheduler.Metrics.RetentionTotal.WithLabelValues(scheduler.Plan.Name, "error").Inc()
}
