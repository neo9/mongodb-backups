package actions

import (
	"runtime"

	"github.com/neo9/mongodb-backups/pkg/api"
	"github.com/neo9/mongodb-backups/pkg/log"
)

func printVersion() {
	log.Info("Go Version: %s", runtime.Version())
	log.Info("Go OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH)
}

func LaunchServer(confPath string, port int32) {
	printVersion()
	backupScheduler := getScheduler(confPath)
	backupScheduler.Run()

	server := &api.HttpServer{
		Port: port,
	}

	log.Info("starting http server on port %v", server.Port)
	server.Start()
}
