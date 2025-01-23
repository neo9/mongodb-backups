package actions

import "github.com/neo9/mongodb-backups/pkg/log"

func DeleteOldBackups(confPath string) {
	log.Info("Triggered delete old backups")
	backupScheduler := getScheduler(confPath)
	backupScheduler.DeleteOldBackups()
}
