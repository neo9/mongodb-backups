package actions

import (
	"os"

	"github.com/neo9/mongodb-backups/pkg/log"
	"github.com/neo9/mongodb-backups/pkg/mongodb"
	"github.com/neo9/mongodb-backups/pkg/scheduler"
	"github.com/neo9/mongodb-backups/pkg/utils"
)

func ArbitraryDump(confPath string) {
	backupScheduler := getScheduler(confPath)
	log.Info("Creating MongoDB dump for %s", backupScheduler.Plan.Name)
	mongoDBDump, err := mongodb.CreateDump(backupScheduler.Plan)
	if err != nil {
		log.Error("Error creating dump for %s", backupScheduler.Plan.Name)
		os.Exit(1)
	}

	UploadDumpFile(mongoDBDump.ArchiveFile, backupScheduler)
	UploadDumpFile(mongoDBDump.LogFile, backupScheduler)
	mongodb.RemoveFile(mongoDBDump.ArchiveFile)
	mongodb.RemoveFile(mongoDBDump.LogFile)
	log.Info("Dump successful")
}

func UploadDumpFile(filename string, scheduler *scheduler.Scheduler) {
	log.Info("Upload file %s. Size: %s", filename, utils.GetHumanFileSize(filename))
	err := scheduler.Bucket.Upload(filename, scheduler.Plan.Name)
	if err != nil {
		log.Error("Could not upload file: %v", err)
		os.Exit(1)
	}
}
