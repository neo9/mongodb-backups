package utils

import "fmt"

func GetHumanBytes(bytes int64) string {
	units := []string{"bytes", "KB", "MB", "GB", "TB"}
	index := 0
	for index = 0; index < len(units); index++ {
		if bytes < 1024 {
			break
		}
		bytes = bytes / 1024
	}

	return fmt.Sprintf("%d %s", bytes, units[index])
}
