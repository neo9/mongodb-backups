package test_integration

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/neo9/mongodb-backups/pkg/utils"
)

const remove_script = "./remove_data.js"
const init_script = "./init.js"
const check_data_script = "./check_data.js"

const max_delay = 20 * time.Second

func remove_data() {
	print("Removing data ....")
	mongosh_command := mongosh_command("test", "test", remove_script)
	output, err := utils.LaunchCommand(mongosh_command, max_delay)

	if err != nil {
		panic(fmt.Sprintf("We couldn't empty the database %s", output))
	}
}

func init_data() {
	print("Init data ....")
	mongosh_command := mongosh_command("test", "test", init_script)
	output, err := utils.LaunchCommand(mongosh_command, max_delay)

	if err != nil {
		panic(fmt.Sprintf("We couldn't init the database %s", string(output)))
	}

}

func number_of_documents() int {
	mongosh_command := mongosh_command("test", "test", check_data_script)
	output, err := utils.LaunchCommand(mongosh_command, max_delay)

	if err != nil {
		panic(fmt.Sprintf("We couldn't init the database %s", string(output)))
	}
	format_output := strings.TrimSpace(string(output))
	number, err := strconv.Atoi(format_output)

	if err != nil {
		panic("No possible conversion to number")
	}

	return number

}

func mongosh_command(user string, password string, script_path string) string {
	return fmt.Sprintf(
		"mongosh -u %s -p %s %s",
		user,
		password,
		script_path,
	)
}
