# MongoDB backup

[![Build Status](https://travis-ci.org/neo9/mongodb-backups.svg?branch=master)](https://travis-ci.org/neo9/mongodb-backups)


Create MongoDB snapshots to an encrypted S3 bucket.
Handle snapshot restoration from backup.
Can be easily monitored by Prometheus.

[Docker repository](https://hub.docker.com/r/neo9sas/mongodb-backups)

## Usage

```bash
# Launch server & scheduler
./mongodb-backups --config ./config.yaml
# List backup
./mongodb-backups --config ./config.yaml --list
# Restore specific backup
./mongodb-backups --config ./config.yaml --restore [id] --args '--drop'
# Restore last backup
./mongodb-backups --config ./config.yaml --restore-last --args '--drop'
```

Parameters:

- `--config`: Config path. Default `./config.yaml`
- `--list`: list backups
- `--restore`: Restore specific backup from snapshot
- `--restore-last`: Restore last backup from snaphost
- `--args`: MongoDB restore additional arguments

## Config file

- `name`: backup name
- `schedule`: cronjob schedule. Example: `0 * * * *`
- `retention`: max retention. Example: `2d`, `1w`, `1M`, `720h`
- `timeout`: mongodb dump timeout
- `mongodb`:
    - `host`: MongoDB host
    - `port`: MongoDB port
- `bucket`:
    - `s3`:
        - `name`: bucket name
        - `region`: bucket region

Example:

```yaml
name: integration
retention: 1w
schedule: '0 0 * * *'
timeout: 15m
mongodb:
  host: localhost
  port: 27017
bucket:
  s3:
    name: bucket-name
    region: eu-west-1
```

## Prometheus metrics

- `mongodb_backups_scheduler_backup_total`: Total number of backups (status: success / error)
- `mongodb_backups_scheduler_retention_total`: Total number of successful retention cleanup (status: success / error)
- `mongodb_backups_scheduler_bucket_snapshot_count`: Current number of snapshots in the bucket
- `mongodb_backups_scheduler_snapshot_size`: Last snapshot size in bytes
- `mongodb_backups_scheduler_snapshot_latency`: Last snapshot duration in seconds

All metrics have the label `name` equals to the config `name` key.

## Environment variables

- `MONGODB_USER`: MongoDB user
- `MONGODB_PASSWORD`: MongoDB password

## AWS

S3 policy example:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "0",
      "Effect": "Allow",
      "Action": [
        "s3:PutObject",
        "s3:GetObject",
        "s3:ListObject",
        "s3:DeleteObject"
      ],
      "Resource": [
        "arn:aws:s3:::my-bucket-name",
        "arn:aws:s3:::my-bucket-name/*"
      ]
    }
  ]
}
```

## Development

### Run

```bash
# With Go
go run ./cmd --config config.yaml

# With Docker
docker build -t n9-backup .
docker run --rm -v /tmp/config:/tmp/config n9-backup mongodb-backup --config /tmp/config/config.yaml
```

