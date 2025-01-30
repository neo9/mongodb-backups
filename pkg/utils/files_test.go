//go:build !integration
// +build !integration

package utils

import "testing"

func TestOkSnapshotLog(t *testing.T) {
	timestamp, _ := GetBucketFileTimestamp("mongodb-snapshot-1631079123456.log")

	if timestamp != 1631079123456 {
		t.Errorf("Expected %d got %d", 1631079123456, &timestamp)
	}
}

func TestOkSnapshotGz(t *testing.T) {
	timestamp, _ := GetBucketFileTimestamp("mongodb-snapshot-1631079123456.gz")

	if timestamp != 1631079123456 {
		t.Errorf("Expected %d got %d", 1631079123456, &timestamp)
	}
}

func TestKOFormat(t *testing.T) {
	_, err := GetBucketFileTimestamp("mongodb-snapshot-22d01m2024y.gz")

	if err == nil {
		t.Errorf("Expected error on the bucket")
	}
}

func TestKOTimestamp(t *testing.T) {
	_, err := GetBucketFileTimestamp("mongodb-snapshot-16310791234563424234324.gz")

	if err == nil {
		t.Errorf("Expected error on the bucket")
	}
}
