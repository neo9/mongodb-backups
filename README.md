# MongoDB backup

Backup MongoDB dumps to S3 or GCS.

**Work in progress**

## TODO

- Retention
- Prometheus metrics and alerting
- GCS support
- Tests
- Better error handler

## Config file

- `name`: backup name
- `schedule`: cronJob schedule. Example: `0 * * * *`
- `mongodb`:
    - `host`: MongoDB host
    - `port`: MongoDB port
- `bucket`: dictionary
    - `s3`:
        - `name`: bucket name
        - `region`: bucket region


## Environment variables

- `MONGO_USER`: MongoDB user
- `MONGO_PASSWORD`: MongoDB password

