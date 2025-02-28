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
# Arbitrary dump
./mongodb-backups --config ./config.yaml --dump
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
- `tmpPath`: path to store tempory backup before s3 upload
- `mongodb`:
    - `host`: MongoDB host
    - `port`: MongoDB port
- `bucket` (fill only one option):
    - `s3`:
        - `name`: bucket name
        - `region`: bucket region
  - `gs`:
      - `name`: bucket name
  - `minio`:
      - `name`: bucket name
      - `host`: bucket hostname (and port if required)
      - `region`: (optional) bucket region
      - `ssl`: (optional) Enable SSL

Example:

```yaml
name: integration
retention: 1w
schedule: '0 0 * * *'
timeout: 15m
tmpPath: /tmp
mongodb:
  host: localhost
  port: 27017
createDump:
  maxRetries: 3
  retryDelay: 60
bucket:
  s3:
    name: bucket-name
    region: eu-west-1
```

Example of instant backup start :

```yaml
kubectl apply -f ./k8s/backup.yaml
# then watch status
watch kubectl -n tools get jobs
```

## Prometheus metrics

- `mongodb_backups_scheduler_backup_total`: Total number of backups (status: success / error)
- `mongodb_backups_scheduler_retention_total`: Total number of successful retention cleanup (status: success / error)
- `mongodb_backups_scheduler_bucket_snapshot_count`: Current number of snapshots in the bucket
- `mongodb_backups_scheduler_snapshot_size`: Last snapshot size in bytes
- `mongodb_backups_scheduler_snapshot_latency`: Last snapshot duration in seconds
- `mongodb_backups_scheduler_last_successful_snaphot`: Last successful snapshot timestamp

All metrics have the label `name` equals to the config `name` key.

## Environment variables

### Global

- `MONGODB_USER`: MongoDB user
- `MONGODB_PASSWORD`: MongoDB password
- `MONGODB_AUTH_ARGS`: MongoDB additional authentication arguments

### AWS

- `AWS_ACCESS_KEY_ID`: AWS secret key ID
- `AWS_SECRET_ACCESS_KEY`: AWS secret access key

### Minio

- `MINIO_ACCESS_KEY_ID`: Minio secret key ID
- `MINIO_SECRET_ACCESS_KEY`: Minio secret access key

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

The changelog is generated by [git-chglog](https://github.com/git-chglog/git-chglog)

You must follow the given commit syntax: `<item>: <description>`

Generate the changelog:
```sh
git-chglog -o CHANGELOG.md
```

### Run

```bash
# With Go
go run ./cmd --config config.yaml

# With Docker
docker build -t n9-backup .
docker run --rm -v /tmp/config:/tmp/config n9-backup mongodb-backups --config /tmp/config/config.yaml

#With Docker compose
docker-compose -f docker-compose.yaml up 
docker exec n9-backup-dev mongodb-backups --config ./home/config.yaml
```

