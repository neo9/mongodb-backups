package mongodb

import (
	"fmt"
	"github.com/codeskyblue/go-sh"
	"github.com/neo9/mongodb-backups/pkg/config"
	"github.com/neo9/mongodb-backups/pkg/utils"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path"
	"time"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
}

type MongoDBDump struct {
    ArchiveFile string
    LogFile string
}

func CreateDump(plan *config.Plan) (MongoDBDump, error) {
	dumpName := getDumpName()
	outputFile := path.Join("/tmp", dumpName)
	mongoDBDump := MongoDBDump{
		ArchiveFile: outputFile + ".gz",
		LogFile: outputFile + ".log",
	}

	authArgs := getAuthenticationArguments()
	dumpCommand := fmt.Sprintf("mongodump %s --archive=%v --gzip --host %s --port %s",
		authArgs,
		mongoDBDump.ArchiveFile,
		plan.MongoDB.Host,
		plan.MongoDB.Port)

	duration, err := utils.GetDurationFromTimeString(plan.Timeout)
	if err != nil {
		log.Errorf("Error parsing timeout: %v", err)
		return mongoDBDump, err
	}

	output, err := sh.Command("/bin/sh", "-c", dumpCommand).
		SetTimeout(duration).
		CombinedOutput()


	if err != nil {
		log.Errorf("Error creating dump: %v", err)
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

func removeFile(filename string) {
	err := os.Remove(filename)
	if err != nil {
		log.Errorf("Cannot delete temporary file %s: %v", filename, err)
	}
}

func getDumpName() string {
	return fmt.Sprintf("mongodb-snapshot-%d", time.Now().Unix())
}

func getAuthenticationArguments() string {
	username, isUsernameDefined := os.LookupEnv("MONGODB_USER")
	password, isPasswordDefined := os.LookupEnv("MONGODB_PASSWORD")

	if isUsernameDefined && isPasswordDefined {
		return fmt.Sprintf("-u %s --password %s", username, password)
	}

	return ""
}

func logToFile(filename string, data []byte) error {
	if len(data) > 0 {
		err := ioutil.WriteFile(filename, data, 0644)
		if err != nil {
			log.Errorf("Error writing dump log file: %v", err)
			return err
		}
	}

	return nil
}
