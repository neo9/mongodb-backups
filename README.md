# MongoDB backup

Backup MongoDB dumps to S3 or GCS.

**Work in progress**

## TODO

- MongoDB authentication support with env variables
- Retention
- Timeout configuration
- Code refactoring
- Prometheus metrics and alerting
- GCS support
- Tests
- Better error handler

## Config file

- `backups`: array
    - `name`: backup name
    - `schedule`: cronJob schedule. Example: `0 * * * *`
    - `mongodb`:
        - `host`: MongoDB host
        - `port`: MongoDB port
- `s3`:
    - `name`: bucket name
    - `region`: bucket region



