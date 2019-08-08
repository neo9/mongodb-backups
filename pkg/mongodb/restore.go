package mongodb

import (
	"fmt"
	"github.com/codeskyblue/go-sh"
	"time"
)

func RestoreDump(filename string, args string, timeout time.Duration) error {
	authArgs := getAuthenticationArguments()
	restoreCommand := fmt.Sprintf(
		"mongorestore --authenticationDatabase admin %s --archive=%s --gzip %s",
		authArgs,
		filename,
		args,
	)
	fmt.Print(restoreCommand)

	err := sh.Command("/bin/sh", "-c", restoreCommand).
		SetTimeout(timeout * time.Second).
		Run()

	return err
}
