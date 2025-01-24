package test_integration

import (
	"fmt"
	"os"
	"sort"
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
	init_data()
	time.Sleep(time.Second)
	actions.DeleteOldBackups(plan)

	//Run the tests
	exitVal := m.Run()

	//Run before testing
	os.Exit(exitVal)
}

func TestArbitraryDump(t *testing.T) {
	actions.ArbitraryDump(plan)
	backups := actions.ListBackups(plan)
	if len(backups) != 1 {
		t.Errorf("The expected number of backups is one")
	}
	time.Sleep(time.Second)
	actions.ArbitraryDump(plan)
	backups = actions.ListBackups(plan)
	if len(backups) != 2 {
		t.Errorf("The expected number of backups is two")
	}
}

func TestRestoreLastBackup(t *testing.T) {
	actions.ArbitraryDump(plan)

	remove_data()
	time.Sleep(time.Second)

	documents := number_of_documents()

	if documents != 0 {
		t.Errorf("Expected %d instead got %d", 0, documents)
	}
	actions.RestoreLastBackup(plan, "--drop")
	time.Sleep(time.Second)
	documents = number_of_documents()

	if documents != NUMBER_OF_POKEMONS {
		t.Errorf("Expected %d instead got %d", NUMBER_OF_POKEMONS, documents)
	}
}

func TestRestoreMediumBackup(t *testing.T) {
	actions.ArbitraryDump(plan)

	remove_data()
	time.Sleep(time.Second)

	actions.ArbitraryDump(plan)

	init_data()
	time.Sleep(time.Second)

	actions.ArbitraryDump(plan)

	backups := actions.ListBackups(plan)
	sort.Slice(backups, func(i, j int) bool {
		return backups[i].Timestamp.Before(backups[j].Timestamp)
	})
	backup_wihtout_data := backups[1]

	if len(backups) != 3 {
		t.Errorf("Expected %d instead got %d", 3, len(backups))
	}

	documents := number_of_documents()
	print(fmt.Sprintf("The length of the test is %d", len(backups)))
	if documents != NUMBER_OF_POKEMONS {
		t.Errorf("Expected %d instead got %d", NUMBER_OF_POKEMONS, documents)
	}

	actions.RestoreBackup(plan, backup_wihtout_data.Etag, "--drop")
	time.Sleep(time.Second)

	documents = number_of_documents()

	if documents != 0 {
		t.Errorf("Expected %d instead got %d", 0, documents)
	}
}
