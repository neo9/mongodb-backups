package config

import (
	"errors"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

const DEFAULT_TRIES = 3
const DEFAULT_DELAY = 60

type Plan struct {
	Name       string     `json:"name"`
	Schedule   string     `json:"schedule"`
	Retention  string     `json:"retention"`
	Timeout    string     `json:"timeout"`
	TmpPath    string     `json:"tmpPath"`
	MongoDB    MongoDB    `json:"mongodb"`
	Bucket     Bucket     `json:"buckets"`
	CreateDump CreateDump `json:"createDump"`
}

type Bucket struct {
	S3    S3    `json:"s3"`
	GS    GS    `json:"gs"`
	Minio Minio `json:"minio"`
}

type S3 struct {
	Name   string `json:"name"`
	Region string `json:"region"`
}
type GS struct {
	Name string `json:"name"`
}

type Minio struct {
	Name   string `json:"name"`
	Host   string `json:"host"`
	Region string `json:"region,omitempty"`
	SSL    bool   `json:"ssl,omitempty"`
}

type CreateDump struct {
	MaxRetries int           `json:"maxRetries"`
	RetryDelay time.Duration `json:"retryDelay"`
}

type MongoDB struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

func (plan *Plan) setDefaults() {
	if plan.CreateDump.MaxRetries == 0 {
		plan.CreateDump.MaxRetries = DEFAULT_TRIES
	}
	if plan.CreateDump.RetryDelay == 0 {
		plan.CreateDump.RetryDelay = DEFAULT_DELAY
	}
}

func (plan *Plan) GetPlan(filename string) (*Plan, error) {
	_, err := os.Stat(filename)
	if err != nil {
		return nil, err
	}

	yamlFile, err := os.ReadFile(filename)
	if err != nil {
		return plan, err
	}

	err = yaml.Unmarshal(yamlFile, plan)
	if err != nil {
		return plan, err
	}
	plan.setDefaults()

	return plan, validate(plan)
}

func validate(plan *Plan) error {
	if plan.Bucket.S3.Name != "" && plan.Bucket.S3.Region != "" && plan.Bucket.GS.Name == "" && plan.Bucket.Minio.Name == "" {
		return nil
	} else if plan.Bucket.GS.Name != "" && plan.Bucket.S3.Name == "" && plan.Bucket.Minio.Name == "" {
		return nil
	} else if plan.Bucket.Minio.Name != "" && plan.Bucket.Minio.Host != "" && plan.Bucket.S3.Name == "" && plan.Bucket.GS.Name == "" {
		return nil
	}

	return errors.New("error in configuration : should only have one type of bucket configured")
}
