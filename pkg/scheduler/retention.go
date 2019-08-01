package scheduler

import (
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
    	log.Errorf("Could not list files for plan %s", scheduler.Plan.Name)
    	return
	}

    retentionDuration, err := utils.GetDurationFromTimeString(scheduler.Plan.Retention)
    if err != nil {
    	log.Errorf("Could not execute retention: %v", err)
    	return
	}

    var removeFiles []string

    for i := 0; i < len(files); i++ {
    	log.Info("File: ", files[i])
    	reg := regexp.MustCompile(`mongodb-snapshot-(?P<Time>\d+)\.(gz|log)`)
    	match := reg.FindStringSubmatch(files[i])
    	if len(match) != 3 {
    		log.Errorf("File does not match pattern in folder %s: %s", scheduler.Plan.Name, files[i])
		}

        timestamp, err := strconv.ParseInt(match[1], 10, 64)
        if err != nil {
			log.Errorf("File has invalid timestamp in folder %s: %s", scheduler.Plan.Name, files[i])
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

    log.Infof("%d file(s) are being removed", len(removeFiles))

    for i := 0; i < len(removeFiles); i++ {
    	err := scheduler.Bucket.DeleteFile(removeFiles[i])
    	if err != nil {
    		log.Errorf("Could not remove file %s", removeFiles[i])
		}
	}
}
