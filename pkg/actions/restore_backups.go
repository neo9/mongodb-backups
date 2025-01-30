package actions

import (
	"os"

	"github.com/neo9/mongodb-backups/pkg/restore"
)

func RestoreBackup(confPath string, restoreID string, args string) {
	backupScheduler := getScheduler(confPath)
	err := restore.Restore(backupScheduler, restoreID, args)
	if err != nil {
		os.Exit(1)
	}
}

func RestoreLastBackup(confPath string, args string) {
	backupScheduler := getScheduler(confPath)
	err := restore.RestoreLast(backupScheduler, args)
	if err != nil {
		os.Exit(1)
	}
}
