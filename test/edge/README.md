# MongoDB Backup and Restore Integration Tests

This folder contains edge case tests for verifying backup and restore functionality under a big load in MongoDB using `mongodb-backups`, this tool uses behind the curtain `mongodump` and `mongorestore`. The tests use a bucket from `upload-minio.sh` the bucket by defualt is the last unknown backup of uber-production since we have the following issue [AONRUN-1248](https://aon-mercure.atlassian.net/browse/AONRUN-1248). 

The environment for testing is powered by Docker, using `docker-compose` to set up the required services.

---

## Prerequisites

Before running the tests, ensure you have the following tools installed:

- **Docker**: To run the MongoDB and MinIO services using `docker-compose -f docker-compose-integration.yaml`.
- **mongodump**: A tool to create MongoDB backups.
- **mongorestore**: A tool to restore MongoDB backups.
- **mongosh**: A MongoDB shell for interacting with the database and launch the commands for testing.

Seems we need the use of Time.Sleep as mongodb does not seem to refresh in such a short time
---

## Run the test

As this test usually takes one hour to do 50 backups I'll recommend to launch it on the bash with

```
go test -timeout 3600s -run ^TestRestoreLastBackup$ github.com/neo9/mongodb-backups/test/edge
```