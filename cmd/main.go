package main

import (
	"flag"
	"fmt"
	"github.com/neo9/mongodb-backups/pkg/api"
	"github.com/neo9/mongodb-backups/pkg/config"
	"github.com/neo9/mongodb-backups/pkg/listing"
	"github.com/neo9/mongodb-backups/pkg/scheduler"
	log "github.com/sirupsen/logrus"
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
	listing.DisplayBackups(backupScheduler)
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

	flag.Parse()
	if *list {
		log.SetFormatter(&log.TextFormatter{})
		listBackups(*confPath)
	} else {
		launchServer(*confPath, int32(*port))
		log.SetFormatter(&log.JSONFormatter{})
	}
}
