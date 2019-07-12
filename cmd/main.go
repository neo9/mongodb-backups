package main


import (
	"flag"
	"fmt"
	"github.com/neo9/mongodb-backups/pkg/api"
	"github.com/neo9/mongodb-backups/pkg/config"
	"github.com/neo9/mongodb-backups/pkg/scheduler"
	log "github.com/sirupsen/logrus"
	"runtime"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
}

func printVersion() {
	log.Info(fmt.Sprintf("Go Version: %s", runtime.Version()))
	log.Info(fmt.Sprintf("Go OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH))
}


func main() {
	printVersion()
	confPath := flag.String("path", "./config.yaml", "Plan config path")
	flag.Parse()
	log.Infof("Parsing configuration file: %s", *confPath)

	plan := config.Plan{}
	_, err := plan.GetPlan(*confPath)
	if err != nil {
		panic(err)
	}

	backupScheduler := scheduler.New(&plan)
	backupScheduler.Run()

	server := &api.HttpServer{
		Port: 8080,
	}
	log.Infof("starting http server on port %v", server.Port)
	server.Start()
}
