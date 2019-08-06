package main


import (
	"flag"
	"fmt"
	"github.com/neo9/mongodb-backups/pkg/api"
	"github.com/neo9/mongodb-backups/pkg/config"
	"github.com/neo9/mongodb-backups/pkg/scheduler"
	log "github.com/sirupsen/logrus"
	"runtime"
	"strconv"
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
	confPath := flag.String("config", "./config.yaml", "Plan config path")
	portStr := flag.String("port", "8080", "Server port")
	flag.Parse()
	log.Infof("Parsing configuration file: %s", *confPath)

	plan := config.Plan{}
	_, err := plan.GetPlan(*confPath)
	if err != nil {
		panic(err)
	}

	port, err := strconv.Atoi(*portStr)
	if err != nil {
		panic(fmt.Sprintf("Invalid server port %s: %v", portStr, err))
	}

	backupScheduler := scheduler.New(&plan)
	backupScheduler.Run()

	server := &api.HttpServer{
		Port: int32(port),
	}
	log.Infof("starting http server on port %v", server.Port)
	server.Start()
}
