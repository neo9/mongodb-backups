package listing

import (
	"github.com/neo9/mongodb-backups/pkg/scheduler"
	"github.com/neo9/mongodb-backups/pkg/utils"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

func DisplayBackups(scheduler *scheduler.Scheduler, ) {
	files, err := scheduler.Bucket.ListFiles(scheduler.Plan.Name)
	if err != nil {
		log.Errorf("Could not list files for %s: %v", scheduler.Plan.Name, err)
	}

	if len(files) == 0 {
		log.Infof("Bucket is empty")
		return
	}

	for i := 0; i < len(files); i++ {
		file := files[i]
		if strings.Contains(file.Name, ".log") {
			continue
		}

		timestamp, err := utils.GetBucketFileTimestamp(file.Name)
		if err != nil {
			log.Errorf("Could not parse file: %v", err)
			continue
		}

		log.Infof("%s | Backup %s, size: %s",
			time.Unix(timestamp, 10),
			file.Etag,
			utils.GetHumanBytes(file.Size))
	}
}
