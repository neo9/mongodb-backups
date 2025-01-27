//go:build edge
// +build edge

package test_edge

import (
	"os"
	"testing"
	"time"

	"github.com/neo9/mongodb-backups/pkg/actions"
)

const plan = "plan.yaml"
const NUMBER_OF_POKEMONS = 151

func TestMain(m *testing.M) {
	//Do before test
	os.Setenv("MINIO_ACCESS_KEY_ID", "minioadmin")
	os.Setenv("MINIO_SECRET_ACCESS_KEY", "minioadmin")
	os.Setenv("MONGODB_USER", "test")
	os.Setenv("MONGODB_PASSWORD", "test")

	remove_data()
	time.Sleep(time.Second)

	//Run the tests
	exitVal := m.Run()

	//Run before testing
	os.Exit(exitVal)
}

func TestRestoreLastBackup(t *testing.T) {
	for i := 0; i < 50; i++ {
		actions.RestoreLastBackup(plan, "--drop")
		actions.ArbitraryDump(plan)
		remove_data()
		time.Sleep(time.Second)
		print("Succeded loop %d", i)
	}

}
