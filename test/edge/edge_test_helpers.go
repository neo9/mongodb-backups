package test_edge

import (
	"fmt"
	"time"

	"github.com/neo9/mongodb-backups/pkg/utils"
)

const remove_script = "./remove_data.js"

const max_delay = 20 * time.Second

func remove_data() {
	print("Removing data ....")
	mongosh_command := mongosh_command("test", "test", remove_script)
	output, err := utils.LaunchCommand(mongosh_command, max_delay)

	if err != nil {
		panic(fmt.Sprintf("We couldn't empty the database %s", output))
	}
}

func mongosh_command(user string, password string, script_path string) string {
	return fmt.Sprintf(
		"mongosh -u %s -p %s %s",
		user,
		password,
		script_path,
	)
}
