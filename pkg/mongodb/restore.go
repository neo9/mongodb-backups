package mongodb

import (
	"fmt"
	"time"

	"github.com/neo9/mongodb-backups/pkg/config"
	"github.com/neo9/mongodb-backups/pkg/log"
	"github.com/neo9/mongodb-backups/pkg/utils"
)

func RestoreDump(filename string, args string, plan *config.Plan) error {
	authArgs := getAuthenticationArguments()
	restoreCommand := fmt.Sprintf(
		"mongorestore --authenticationDatabase admin %s --archive=%s --gzip %s --host %s --port %s",
		authArgs,
		filename,
		args,
		plan.MongoDB.Host,
		plan.MongoDB.Port)
	fmt.Print(restoreCommand)

	duration, err := utils.GetDurationFromTimeString(plan.Timeout)
	if err != nil {
		log.Error("Error parsing timeout: %v", err)
		return err
	}

	startTime := time.Now()
	output, err := utils.LaunchCommand(restoreCommand, duration)
	seconds := time.Since(startTime).Seconds()

	if err != nil {
		log.Error("Error restoring dump: %v, %s", err, output)
		log.Error("Restoring timeout: %s", duration)
		return err
	}

	log.Info("Done restoring dump for %s in %f s", plan.Name, seconds)

	return err
}
