package mongodb

import (
	"fmt"
	"github.com/codeskyblue/go-sh"
	"github.com/neo9/mongodb-backups/pkg/config"
	"github.com/neo9/mongodb-backups/pkg/utils"
	"github.com/prometheus/common/log"
	"path"
	"time"
)

type MongoDBDump struct {
	ArchiveFile string
	LogFile string
	Duration float64
}

func CreateDump(plan *config.Plan) (MongoDBDump, error) {
	dumpName := getDumpName()
	outputFile := path.Join("/tmp", dumpName)
	mongoDBDump := MongoDBDump{
		ArchiveFile: outputFile + ".gz",
		LogFile: outputFile + ".log",
	}

	authArgs := getAuthenticationArguments()
	dumpCommand := fmt.Sprintf(
		"mongodump --authenticationDatabase admin %s --archive=%v --gzip --host %s --port %s",
		authArgs,
		mongoDBDump.ArchiveFile,
		plan.MongoDB.Host,
		plan.MongoDB.Port)

	duration, err := utils.GetDurationFromTimeString(plan.Timeout)
	if err != nil {
		log.Errorf("Error parsing timeout: %v", err)
		return mongoDBDump, err
	}

	startTime := time.Now()
	output, err := sh.Command("/bin/sh", "-c", dumpCommand).
		SetTimeout(duration).
		CombinedOutput()
	mongoDBDump.Duration = time.Since(startTime).Seconds()


	if err != nil {
		log.Errorf("Error creating dump: %v, %s", err, output)
		log.Errorf("Dump timeout: %s", duration)
		removeFile(mongoDBDump.ArchiveFile)
		return mongoDBDump, err
	}

	log.Infof("Done creating dump for %s", plan.Name)
	err = logToFile(mongoDBDump.LogFile, output)
	if err != nil {
		removeFile(mongoDBDump.ArchiveFile)
	}

	return mongoDBDump, err
}

func getDumpName() string {
	return fmt.Sprintf("mongodb-snapshot-%d", time.Now().Unix())
}

