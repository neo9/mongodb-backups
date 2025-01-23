package actions

import (
	"github.com/neo9/mongodb-backups/pkg/config"
	"github.com/neo9/mongodb-backups/pkg/log"
	"github.com/neo9/mongodb-backups/pkg/scheduler"
)

func getScheduler(confPath string) *scheduler.Scheduler {
	log.Info("Parsing configuration file: %s", confPath)
	plan := config.Plan{}
	_, err := plan.GetPlan(confPath)
	if err != nil {
		panic(err)
	}

	return scheduler.New(&plan)
}
