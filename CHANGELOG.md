
<a name="v0.12.0"></a>
## [v0.12.0](https://github.com/neo9/mongodb-backups/compare/v0.11.0...v0.12.0) (2025-01-02)

### Mongodump

* add retries in CreateDump + configuration to set maxRetries and delay time

### Release

* v0.12.0


<a name="v0.11.0"></a>
## [v0.11.0](https://github.com/neo9/mongodb-backups/compare/v0.10.0...v0.11.0) (2021-10-14)

### Mongo

* use mongodb-tools 100.5.0

### Release

* v0.11.0


<a name="v0.10.0"></a>
## [v0.10.0](https://github.com/neo9/mongodb-backups/compare/v0.9.0...v0.10.0) (2021-10-11)

### Actions

* Fix target branch

### Mongo

* Define both possibilities
* Add custom authentication args

### Release

* v0.10.0
* Add release-it configuration


<a name="v0.9.0"></a>
## [v0.9.0](https://github.com/neo9/mongodb-backups/compare/v0.8.0...v0.9.0) (2021-09-20)

### Bucket

* Add minio support

### Config

* Fix bucket configuration check

### Doc

* Add dockerfile fix in changelog
* Update changelog
* Add changelog specifications
* Add minio documentation
* Add changelog

### Fix

* Upgrade golang version in dockerfile

### Fmt

* Fix syntax


<a name="v0.8.0"></a>
## [v0.8.0](https://github.com/neo9/mongodb-backups/compare/v0.7.0...v0.8.0) (2021-07-30)

### Config

* Add variable to choose tmpPath


<a name="v0.7.0"></a>
## [v0.7.0](https://github.com/neo9/mongodb-backups/compare/v0.6.1...v0.7.0) (2020-12-18)

### Ci

* Move on github actions

### Dump

* Fix removefile after dump

### Gs

* Add gs bucket support

### Job

* Add instant backup example

### Removefile

* Make this function exportable

### Restore

* Use configured timeout on restore


<a name="v0.6.1"></a>
## [v0.6.1](https://github.com/neo9/mongodb-backups/compare/v0.6.0...v0.6.1) (2020-08-26)

### Dump

* Remove tmp file after backup

### K8s

* Increase memory limit and request to 1Gi

### Prom

* Add alert examples


<a name="v0.6.0"></a>
## [v0.6.0](https://github.com/neo9/mongodb-backups/compare/v0.5.1...v0.6.0) (2019-08-09)

### Doc

* Add new prometheus metrics

### Pkg

* Add last successful backup timestamp


<a name="v0.5.1"></a>
## [v0.5.1](https://github.com/neo9/mongodb-backups/compare/v0.5.0...v0.5.1) (2019-08-09)

### Travis

* Set file glob to true
* Add darwin build support


<a name="v0.5.0"></a>
## [v0.5.0](https://github.com/neo9/mongodb-backups/compare/v0.4.1...v0.5.0) (2019-08-09)

### Pkg

* Add --dump support


<a name="v0.4.1"></a>
## [v0.4.1](https://github.com/neo9/mongodb-backups/compare/v0.4.0...v0.4.1) (2019-08-09)

### Doc

* Add config example
* Add S3 policy example

### K8s

* Improve deployment specs

### Pkg

* Fix metrics inconsistent label cardinality


<a name="v0.4.0"></a>
## [v0.4.0](https://github.com/neo9/mongodb-backups/compare/v0.3.0...v0.4.0) (2019-08-09)

### Pkg

* Fix current metrics types


<a name="v0.3.0"></a>
## [v0.3.0](https://github.com/neo9/mongodb-backups/compare/v0.2.0...v0.3.0) (2019-08-08)

### Pkg

* Add restore command


<a name="v0.2.0"></a>
## [v0.2.0](https://github.com/neo9/mongodb-backups/compare/v0.1.0...v0.2.0) (2019-08-06)

### Cmd

* Add port override


<a name="v0.1.0"></a>
## v0.1.0 (2019-08-02)

### Ci

* Fix docker tag
* Fix travis
* Add travis

### Doc

* Fix image names
* Add travis badge
* Add docker
* Add README

### Docker

* Fix command line name
* Add ca-certificates package
* Add Dockerfile

### Metrics

* Add basic support

### Mongodb

* Add better file cleanup

### Pkg

* Add authentication database
* Add log timeout error
* Add retention metrics
* Fix MongoDB size metrics
* Add metrics duration
* Add MongoDB logs output
* Fix metrics error
* Add retention support

### Refactor

* Add new backup logic
* Types and bucket logic

### Travis

* Add build command

