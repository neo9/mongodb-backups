package mongodb

import (
	"fmt"
	"path"
	"time"

	"github.com/neo9/mongodb-backups/pkg/config"
	"github.com/neo9/mongodb-backups/pkg/utils"
)

type MongoDBDump struct {
	ArchiveFile string
	LogFile     string
	Duration    float64
}

func CreateDump(plan *config.Plan) (MongoDBDump, error) {
	var err error
	var mongoDBDump MongoDBDump

	maxRetries := plan.CreateDump.MaxRetries
	retryDelay := plan.CreateDump.RetryDelay * time.Second

	for i := 0; i < maxRetries; i++ {
		mongoDBDump, err = CreateDumpInternal(plan)
		if err != nil {
			utils.Error("Error creating mongodump (retry %d/%d): %v", i+1, maxRetries, err)
			time.Sleep(retryDelay)
		} else {
			return mongoDBDump, nil
		}
	}

	return mongoDBDump, err
}

func CreateDumpInternal(plan *config.Plan) (MongoDBDump, error) {
	dumpName := getDumpName()
	outputFile := path.Join(plan.TmpPath, dumpName)
	mongoDBDump := MongoDBDump{
		ArchiveFile: outputFile + ".gz",
		LogFile:     outputFile + ".log",
	}

	authArgs := getAuthenticationArguments()
	dumpCommand := fmt.Sprintf(
		"mongodump --forceTableScan --authenticationDatabase admin %s --archive=%v --gzip --host %s --port %s",
		authArgs,
		mongoDBDump.ArchiveFile,
		plan.MongoDB.Host,
		plan.MongoDB.Port)

	duration, err := utils.GetDurationFromTimeString(plan.Timeout)
	if err != nil {
		utils.Error("Error parsing timeout: %v", err)
		return mongoDBDump, err
	}

	startTime := time.Now()
	output, err := utils.LaunchCommand(dumpCommand, duration)
	mongoDBDump.Duration = time.Since(startTime).Seconds()

	if err != nil {
		utils.Error("Error creating dump: %v, %s", err, output)
		utils.Error("Dump timeout: %s", duration)
		RemoveFile(mongoDBDump.ArchiveFile)
		return mongoDBDump, err
	}

	utils.Info("Done creating dump for %s", plan.Name)
	err = logToFile(mongoDBDump.LogFile, output)
	if err != nil {
		RemoveFile(mongoDBDump.ArchiveFile)
	}

	return mongoDBDump, err
}

func getDumpName() string {
	return fmt.Sprintf("mongodb-snapshot-%d", time.Now().Unix())
}
