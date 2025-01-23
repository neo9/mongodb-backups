# MongoDB Backup and Restore Integration Tests

This repository contains integration tests for verifying backup and restore functionality in MongoDB using `mongodb-backups`, this tool uses behind the curtain `mongodump` and `mongorestore`. The tests utilize a `seed.json` dataset from the [pokemon repository](https://github.com/ATL-WDI-Exercises/mongo-pokemon), which contains a dataset of trivial set of data. 

The environment for testing is powered by Docker, using `docker-compose` to set up the required services.

---

## Prerequisites

Before running the tests, ensure you have the following tools installed:

- **Docker**: To run the MongoDB and MinIO services using `docker-compose -f docker-compose-integration.yaml`.
- **mongodump**: A tool to create MongoDB backups.
- **mongorestore**: A tool to restore MongoDB backups.
- **mongosh**: A MongoDB shell for interacting with the database and launch the commands for testing.

---
