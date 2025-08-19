package mongodb

import (
        "fmt"
        "os"
        "path"
        "strings"
        "time"

        "github.com/neo9/mongodb-backups/pkg/config"
        "github.com/neo9/mongodb-backups/pkg/log"
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

	mongoDBDump, err = CreateDumpInternal(plan)

	return mongoDBDump, err
}

func CreateDumpInternal(plan *config.Plan) (MongoDBDump, error) {
	dumpName := getDumpName()
	outputFile := path.Join(plan.TmpPath, dumpName)
	mongoDBDump := MongoDBDump{
		ArchiveFile: outputFile + ".gz",
		LogFile:     outputFile + ".log",
	}

        dumpCommand := buildDumpCommand(plan, mongoDBDump.ArchiveFile)

	duration, err := utils.GetDurationFromTimeString(plan.Timeout)
	if err != nil {
		log.Error("Error parsing timeout: %v", err)
		return mongoDBDump, err
	}

	startTime := time.Now()
	output, err := utils.LaunchCommand(dumpCommand, duration)
	mongoDBDump.Duration = time.Since(startTime).Seconds()

	if err != nil {
		log.Error("Error creating dump: %v, %s", err, output)
		log.Error("%s", string(output))
		RemoveFile(mongoDBDump.ArchiveFile)
		return mongoDBDump, err
	}

	log.Info("Done creating dump for %s", plan.Name)
	err = logToFile(mongoDBDump.LogFile, output)
	if err != nil {
		RemoveFile(mongoDBDump.ArchiveFile)
	}

	return mongoDBDump, err
}

func getDumpName() string {
        return fmt.Sprintf("mongodb-snapshot-%d", time.Now().Unix())
}

func buildDumpCommand(plan *config.Plan, archiveFile string) string {
        authArgs := strings.TrimSpace(getAuthenticationArguments())
        if uri, ok := os.LookupEnv("MONGO_URI"); ok {
                cmd := fmt.Sprintf("mongodump --forceTableScan --archive=%v --gzip --uri %s", archiveFile, uri)
                if authArgs != "" {
                        cmd += " " + authArgs
                }
                return cmd
        }
        if authArgs != "" {
                authArgs = " " + authArgs
        }
        return fmt.Sprintf("mongodump --forceTableScan --authenticationDatabase admin%s --archive=%v --gzip --host %s --port %s",
                authArgs,
                archiveFile,
                plan.MongoDB.Host,
                plan.MongoDB.Port)
}
