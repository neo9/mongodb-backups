package utils

import (
	"testing"
	"time"
)

func TestGetMonth(t *testing.T) {
	checkMinutes(t, "2M", 87600)
}

func TestGetWeek(t *testing.T) {
	checkMinutes(t, "3w", 30240)
}

func TestGetMinutes(t *testing.T) {
	checkMinutes(t, "60m", 60)
}

func checkMinutes(t *testing.T, timeStr string, expectedMinutes int) {
	duration, err := GetDurationFromTimeString(timeStr)
	expectedDuration := time.Duration(expectedMinutes) * time.Minute
	if err != nil || duration != expectedDuration {
		t.Errorf("Expected: %d, got %d", expectedMinutes, expectedDuration)
	}
}
