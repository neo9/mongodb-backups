package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
}

func GetDurationFromTimeString(timeStr string) (time.Duration, error) {
	reg := regexp.MustCompile(`(?P<Number>\d+)(?P<Unit>[Mwdhm])`)
	match := reg.FindStringSubmatch(timeStr)
	if len(match) != 3 {
		return 0, fmt.Errorf("Could not parse string: " + timeStr + ". Wrong time format. Example: 1h, 3w, 15d")
	}

	number, _ := strconv.Atoi(match[1])
	unit := match[2]

	unitMap := map[string]int{
		"M": 43800,
		"w": 10080,
		"d": 1440,
		"h": 60,
		"m": 1,
	}

	duration := time.Duration(number*unitMap[unit]) * time.Minute

	return duration, nil
}
