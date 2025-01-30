package main

import (
	"flag"

	"github.com/neo9/mongodb-backups/pkg/actions"
)

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
		actions.ListBackups(*confPath)
	} else if *dump {
		actions.ArbitraryDump(*confPath)
	} else if *restoreLast {
		actions.RestoreLastBackup(*confPath, *args)
	} else if *restoreID != "" {
		actions.RestoreBackup(*confPath, *restoreID, *args)
	} else {
		actions.LaunchServer(*confPath, int32(*port))
	}
}
