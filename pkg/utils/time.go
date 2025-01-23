package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

func GetDurationFromTimeString(timeStr string) (time.Duration, error) {
	reg := regexp.MustCompile(`(?P<Number>\d+)(?P<Unit>[Mwdhms])`)
	match := reg.FindStringSubmatch(timeStr)
	if len(match) != 3 {
		return 0, fmt.Errorf("could not parse string: %s. wrong time format. Example: 1h, 3w, 15d", timeStr)
	}

	number, _ := strconv.Atoi(match[1])
	unit := match[2]

	unitMap := map[string]int{
		"M": 2628000,
		"w": 604800,
		"d": 86400,
		"h": 3600,
		"m": 60,
		"s": 1,
	}

	duration := time.Duration(number*unitMap[unit]) * time.Second

	return duration, nil
}
