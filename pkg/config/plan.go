package config

import (
	"errors"
	"os"
	"time"
	"fmt"

	"gopkg.in/yaml.v2"
)

const DEFAULT_TRIES = 3
const DEFAULT_DELAY = 60

type Plan struct {
	Name       string     `yaml:"name"`
	Schedule   string     `yaml:"schedule"`
	Retention  string     `yaml:"retention"`
	Timeout    string     `yaml:"timeout"`
	TmpPath    string     `yaml:"tmpPath"`
	MongoDB    MongoDB    `yaml:"mongodb"`
	Bucket     Bucket     `yaml:"bucket"`
	CreateDump CreateDump `yaml:"createDump"`
}

type Bucket struct {
	S3    S3    `yaml:"s3"`
	GS    GS    `yaml:"gs"`
	Minio Minio `yaml:"minio"`
}

type S3 struct {
	Name   string `yaml:"name"`
	Region string `yaml:"region"`
}
type GS struct {
	Name string `yaml:"name"`
}

type Minio struct {
	Name   string `yaml:"name"`
	Host   string `yaml:"host"`
	Region string `yaml:"region,omitempty"`
	SSL    bool   `yaml:"ssl,omitempty"`
}

type CreateDump struct {
    MaxRetries int           `yaml:"maxRetries"`
    RetryDelay time.Duration `yaml:"retryDelay"`
}

type MongoDB struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
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


	// 🔍 Verificar valores después de parsear YAML
	fmt.Printf("Parsed YAML - MaxRetries: %d, RetryDelay: %v\n", plan.CreateDump.MaxRetries, plan.CreateDump.RetryDelay)

	plan.setDefaults()

	// 🔍 Verificar valores después de aplicar defaults
	fmt.Printf("After Defaults - MaxRetries: %d, RetryDelay: %v\n", plan.CreateDump.MaxRetries, plan.CreateDump.RetryDelay)

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


