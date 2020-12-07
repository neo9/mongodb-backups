package main

import (
	"flag"
	"fmt"
	"github.com/neo9/mongodb-backups/pkg/api"
	"github.com/neo9/mongodb-backups/pkg/config"
	"github.com/neo9/mongodb-backups/pkg/mongodb"
	"github.com/neo9/mongodb-backups/pkg/restore"
	"github.com/neo9/mongodb-backups/pkg/scheduler"
	"github.com/neo9/mongodb-backups/pkg/utils"
	log "github.com/sirupsen/logrus"
	"os"
	"runtime"
)

func printVersion() {
	log.Info(fmt.Sprintf("Go Version: %s", runtime.Version()))
	log.Info(fmt.Sprintf("Go OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH))
}

func getScheduler(confPath string) *scheduler.Scheduler {
	log.Infof("Parsing configuration file: %s", confPath)
	plan := config.Plan{}
	_, err := plan.GetPlan(confPath)
	if err != nil {
		panic(err)
	}

	return scheduler.New(&plan)
}

func listBackups(confPath string) {
	backupScheduler := getScheduler(confPath)
	err := restore.DisplayBackups(backupScheduler)
	if err != nil {
		os.Exit(1)
	}
}

func restoreBackup(confPath string, restoreID string, args string) {
	backupScheduler := getScheduler(confPath)
	err := restore.Restore(backupScheduler, restoreID, args)
	if err != nil {
		os.Exit(1)
	}
}

func restoreLastBackup(confPath string, args string) {
	backupScheduler := getScheduler(confPath)
	err := restore.RestoreLast(backupScheduler, args)
	if err != nil {
		os.Exit(1)
	}
}

func arbitraryDump(confPath string) {
	backupScheduler := getScheduler(confPath)
	log.Infof("Creating MongoDB dump for %s", backupScheduler.Plan.Name)
	mongoDBDump, err := mongodb.CreateDump(backupScheduler.Plan)
	if err != nil {
		log.Errorf("Error creating dump for %s", backupScheduler.Plan.Name)
		os.Exit(1)
	}

	uploadDumpFile(mongoDBDump.ArchiveFile, backupScheduler)
	uploadDumpFile(mongoDBDump.LogFile, backupScheduler)
	removeFile(mongoDBDump.ArchiveFile)
	removeFile(mongoDBDump.LogFile)
	log.Infof("Dump successful")
}

func uploadDumpFile(filename string, scheduler *scheduler.Scheduler) {
	log.Infof("Upload file %s. Size: %s", filename, utils.GetHumanFileSize(filename))
	err := scheduler.Bucket.Upload(filename, scheduler.Plan.Name)
	if err != nil {
		log.Errorf("Could not upload log file: %v", err)
		os.Exit(1)
	}
}

func launchServer(confPath string, port int32) {
	printVersion()
	backupScheduler := getScheduler(confPath)
	backupScheduler.Run()

	server := &api.HttpServer{
		Port: port,
	}

	log.Infof("starting http server on port %v", server.Port)
	server.Start()
}

func main() {
	confPath := flag.String("config", "./config.yaml", "Plan config path")
	port := flag.Int("port", 8080, "Server port")
	list := flag.Bool("list", false, "List backups")
	dump := flag.Bool("dump", false, "Arbitrary dump")
	restoreID := flag.String("restore", "", "Restore specific backup")
	restoreLast := flag.Bool("restore-last", false, "Restore last backup")
	args := flag.String("args", "", "MongoDB args")

	flag.Parse()
	if *list {
		log.SetFormatter(&log.TextFormatter{})
		listBackups(*confPath)
	} else if *dump {
		log.SetFormatter(&log.TextFormatter{})
		arbitraryDump(*confPath)
	} else if *restoreLast {
		log.SetFormatter(&log.TextFormatter{})
		restoreLastBackup(*confPath, *args)
	} else if *restoreID != "" {
		log.SetFormatter(&log.TextFormatter{})
		restoreBackup(*confPath, *restoreID, *args)
	} else {
		launchServer(*confPath, int32(*port))
		log.SetFormatter(&log.JSONFormatter{})
	}
}

