package test_integration

import (
	"os"
	"testing"

	"github.com/neo9/mongodb-backups/pkg/actions"
)

func TestMain(m *testing.M) {
	os.Setenv("MINIO_ACCESS_KEY_ID", "minioadmin")
	os.Setenv("MINIO_SECRET_ACCESS_KEY", "minioadmin")
	os.Setenv("MONGODB_USER", "test")
	os.Setenv("MONGODB_PASSWORD", "test")
}

func TestArbitraryDump(t *testing.T) {
	actions.ArbitraryDump("plan.yaml")
}

func TestListBackup(t *testing.T) {
	backups := actions.ListBackups("plan.yaml")

	print(backups)
	if len(backups) != 2 {
		t.Errorf("The expected number of backups is one")
	}
}

func TestRestoreBackup(t *testing.T) {
	actions.RestoreLastBackup("plan.yaml", "--drop")
}
