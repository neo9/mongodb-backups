package mongodb

import (
	"fmt"
	"os"

	"github.com/neo9/mongodb-backups/pkg/log"
)

func RemoveFile(filename string) {
	err := os.Remove(filename)
	if err != nil {
		log.Error("Cannot delete temporary file %s: %v", filename, err)
	}
}

func getAuthenticationArguments() string {
	username, isUsernameDefined := os.LookupEnv("MONGODB_USER")
	password, isPasswordDefined := os.LookupEnv("MONGODB_PASSWORD")
	auth_args, isAuthArgsDefined := os.LookupEnv("MONGODB_AUTH_ARGS")
	result := ""

	if isUsernameDefined && isPasswordDefined {
		result = result + fmt.Sprintf("-u %s --password %s", username, password)
	}
	if isAuthArgsDefined {
		result = result + " " + auth_args
	}

	return result
}

func logToFile(filename string, data []byte) error {
	if len(data) > 0 {
		err := os.WriteFile(filename, data, 0644)
		if err != nil {
			log.Error("Error writing dump log file: %v", err)
			return err
		}
	}

	return nil
}
