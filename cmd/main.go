package main

import (
	"flag"
	"fmt"
	"github.com/neo9/mongodb-backups/pkg/api"
	"github.com/neo9/mongodb-backups/pkg/config"
	"github.com/neo9/mongodb-backups/pkg/restore"
	"github.com/neo9/mongodb-backups/pkg/scheduler"
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

func restoreBackup(confPath string, restoreID string) {
	backupScheduler := getScheduler(confPath)
	err := restore.Restore(backupScheduler, restoreID)
	if err != nil {
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
	restoreID := flag.String("restore", "", "Restore specific backup")

	flag.Parse()
	if *list {
		log.SetFormatter(&log.TextFormatter{})
		listBackups(*confPath)
	} else if *restoreID != "" {
		log.SetFormatter(&log.TextFormatter{})
		restoreBackup(*confPath, *restoreID)
	} else {
		launchServer(*confPath, int32(*port))
		log.SetFormatter(&log.JSONFormatter{})
	}
}

