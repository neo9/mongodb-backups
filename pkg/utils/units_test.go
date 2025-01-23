package utils

import "testing"

const TB int64 = 1099511627776
const MB int64 = 1048576
const GB int64 = 1073741824
const KB int64 = 1024

func TestHumanTera(t *testing.T) {
	byteSize := TB + 2*GB + 3*MB + 4*KB + 5
	result := GetHumanBytes(byteSize)

	if result != "1 TB" {
		t.Errorf("Expected %s got %s", result, "1 TB")
	}
}

func TestHumanGB(t *testing.T) {
	byteSize := 2*GB + 3*MB + 4*KB + 5
	result := GetHumanBytes(byteSize)

	if result != "2 GB" {
		t.Errorf("Expected %s got %s", result, "2 GB")
	}
}

func TestHumanMB(t *testing.T) {
	byteSize := 3*MB + 4*KB + 5
	result := GetHumanBytes(byteSize)

	if result != "3 MB" {
		t.Errorf("Expected %s got %s", result, "3 MB")
	}
}

func TestHumanKB(t *testing.T) {
	byteSize := 4*KB + 5
	result := GetHumanBytes(byteSize)

	if result != "4 KB" {
		t.Errorf("Expected %s got %s", result, "4 KB")
	}
}
