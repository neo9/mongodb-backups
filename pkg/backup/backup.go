package backup

import (
	"fmt"
	"github.com/codeskyblue/go-sh"
	"github.com/neo9/mongodb-backups/pkg/config"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"path"
	"time"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
}

func MongoDBDump(mongodb *config.MongoDB, name string) ([]string, error) {
	outputFile := path.Join("/tmp", fmt.Sprintf("mongodb-backup-%s-%d", name, time.Now().Unix()))
	archiveFile := outputFile + ".gz"
	logFile := outputFile + ".log"

	dumpCommand := fmt.Sprintf("mongodump --archive=%v --gzip --host %s --port %s",
		archiveFile,
		mongodb.Host,
		mongodb.Port)

	// TODO: dynamic timeout
	output, err := sh.Command("/bin/sh", "-c", dumpCommand).SetTimeout(time.Duration(5) * time.Minute).CombinedOutput()
	if err != nil {
		log.Errorf("Error creating dump: %v", err)
		return []string{}, err
	}

	log.Infof("Done creating dump for %s", name)

	err = logToFile(logFile, output)
	if err != nil {
		log.Errorf("Error writing dump log file: %v", err)
		return []string{archiveFile}, err
	}

	log.Infof("Done creating log file for %s", name)

	return []string{archiveFile, logFile}, nil
}

func logToFile(filename string, data []byte) error {
	if len(data) > 0 {
		err := ioutil.WriteFile(filename, data, 0644)
		if err != nil {
			return err
		}
	}

	return nil
}