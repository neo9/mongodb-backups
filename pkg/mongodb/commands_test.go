package mongodb

import (
    "os"
    "strings"
    "testing"

    "github.com/neo9/mongodb-backups/pkg/config"
)

func TestBuildDumpCommandWithURI(t *testing.T) {
    os.Unsetenv("MONGODB_USER")
    os.Unsetenv("MONGODB_PASSWORD")
    os.Unsetenv("MONGODB_AUTH_ARGS")
    os.Setenv("MONGO_URI", "mongodb://user:pass@localhost:27017")
    defer os.Unsetenv("MONGO_URI")

    plan := &config.Plan{MongoDB: config.MongoDB{Host: "localhost", Port: "27017"}}
    cmd := buildDumpCommand(plan, "dump.gz")
    if !strings.Contains(cmd, "--uri mongodb://user:pass@localhost:27017") {
        t.Fatalf("expected uri in command, got %s", cmd)
    }
    if strings.Contains(cmd, "--host") || strings.Contains(cmd, "--port") {
        t.Fatalf("did not expect host/port in command: %s", cmd)
    }
}

func TestBuildDumpCommandWithoutURI(t *testing.T) {
    os.Unsetenv("MONGO_URI")
    os.Unsetenv("MONGODB_USER")
    os.Unsetenv("MONGODB_PASSWORD")
    os.Unsetenv("MONGODB_AUTH_ARGS")

    plan := &config.Plan{MongoDB: config.MongoDB{Host: "localhost", Port: "27017"}}
    cmd := buildDumpCommand(plan, "dump.gz")
    if !strings.Contains(cmd, "--host localhost") || !strings.Contains(cmd, "--port 27017") {
        t.Fatalf("expected host and port in command: %s", cmd)
    }
    if strings.Contains(cmd, "--uri") {
        t.Fatalf("did not expect uri in command: %s", cmd)
    }
}

func TestBuildRestoreCommandWithURI(t *testing.T) {
    os.Unsetenv("MONGODB_USER")
    os.Unsetenv("MONGODB_PASSWORD")
    os.Unsetenv("MONGODB_AUTH_ARGS")
    os.Setenv("MONGO_URI", "mongodb://user:pass@localhost:27017")
    defer os.Unsetenv("MONGO_URI")

    plan := &config.Plan{MongoDB: config.MongoDB{Host: "localhost", Port: "27017"}}
    cmd := buildRestoreCommand("dump.gz", "--drop", plan)
    if !strings.Contains(cmd, "--uri mongodb://user:pass@localhost:27017") {
        t.Fatalf("expected uri in command, got %s", cmd)
    }
    if strings.Contains(cmd, "--host") || strings.Contains(cmd, "--port") {
        t.Fatalf("did not expect host/port in command: %s", cmd)
    }
}

func TestBuildRestoreCommandWithoutURI(t *testing.T) {
    os.Unsetenv("MONGO_URI")
    os.Unsetenv("MONGODB_USER")
    os.Unsetenv("MONGODB_PASSWORD")
    os.Unsetenv("MONGODB_AUTH_ARGS")

    plan := &config.Plan{MongoDB: config.MongoDB{Host: "localhost", Port: "27017"}}
    cmd := buildRestoreCommand("dump.gz", "--drop", plan)
    if !strings.Contains(cmd, "--host localhost") || !strings.Contains(cmd, "--port 27017") {
        t.Fatalf("expected host and port in command: %s", cmd)
    }
    if strings.Contains(cmd, "--uri") {
        t.Fatalf("did not expect uri in command: %s", cmd)
    }
}

