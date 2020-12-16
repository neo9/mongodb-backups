package mongodb

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
}


func RemoveFile(filename string) {
	err := os.Remove(filename)
	if err != nil {
		log.Errorf("Cannot delete temporary file %s: %v", filename, err)
	}
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
