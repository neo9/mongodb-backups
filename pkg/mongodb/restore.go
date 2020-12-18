package mongodb

import (
	"fmt"
	"github.com/codeskyblue/go-sh"
	"github.com/neo9/mongodb-backups/pkg/config"
	"github.com/neo9/mongodb-backups/pkg/utils"
	"github.com/prometheus/common/log"
	"time"
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
		log.Errorf("Error parsing timeout: %v", err)
		return err
	}

	startTime := time.Now()
	output, err := sh.Command("/bin/sh", "-c", restoreCommand).
		SetTimeout(duration).
		CombinedOutput()
	seconds := time.Since(startTime).Seconds()

	if err != nil {
		log.Errorf("Error restoring dump: %v, %s", err, output)
		log.Errorf("Restoring timeout: %s", duration)
		return err
	}

	log.Infof("Done restoring dump for %s in %s s", plan.Name, seconds)

	return err
}
