package actions

import (
	"os"

	"github.com/neo9/mongodb-backups/pkg/log"
	"github.com/neo9/mongodb-backups/pkg/restore"
)

func ListBackups(confPath string) []restore.Backup {
	backupScheduler := getScheduler(confPath)
	backups, err := restore.GetBackups(backupScheduler)
	if err != nil {
		os.Exit(1)
	}
	displayBackups(backups)

	return backups
}

func displayBackups(backups []restore.Backup) {
	for _, backup := range backups {
		log.Info("Backup %s | Etag: %s | Size: %s",
			backup.Timestamp,
			backup.Etag,
			backup.Size,
		)
	}
}
