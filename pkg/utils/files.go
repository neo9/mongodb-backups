package utils

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

func GetBucketFileTimestamp(file string) (int64, error) {
	reg := regexp.MustCompile(`mongodb-snapshot-(?P<Time>\d+)\.(gz|log)`)
	match := reg.FindStringSubmatch(file)
	if len(match) != 3 {
		return 0, errors.New(fmt.Sprintf("File does not match pattern in folder %s", file))
	}

	timestamp, err := strconv.ParseInt(match[1], 10, 64)
	if err != nil {
		return 0, errors.New(fmt.Sprintf("File has invalid timestamp in folder: %s", file))
	}

	return timestamp, nil
}

func GetHumanFileSize(filename string) string {
	stat, err := os.Stat(filename)
	if err != nil {
		return "UNKNOWN"
	}
	return GetHumanBytes(stat.Size())
}
