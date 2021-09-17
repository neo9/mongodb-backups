package config

import (
	"errors"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type Plan struct {
	Name      string  `json:"name"`
	Schedule  string  `json:"schedule"`
	Retention string  `json:"retention"`
	Timeout   string  `json:"timeout"`
	TmpPath   string  `json:"tmpPath"`
	MongoDB   MongoDB `json:"mongodb"`
	Bucket    Bucket  `json:"buckets"`
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

type MongoDB struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

func (plan *Plan) GetPlan(filename string) (*Plan, error) {
	_, err := os.Stat(filename)
	if err != nil {
		return nil, err
	}

	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return plan, err
	}

	err = yaml.Unmarshal(yamlFile, plan)
	if err != nil {
		return plan, err
	}

	return plan, validate(plan)
}

func validate(plan *Plan) error {
	if plan.Bucket.S3.Name != "" && plan.Bucket.GS.Name != "" && plan.Bucket.Minio.Name != "" {
		return errors.New("error in configuration : should only have S3, GS or Minio bucket configured")
	}

	if plan.Bucket.S3.Name != "" && plan.Bucket.S3.Region != "" {
		return nil
	}
	if plan.Bucket.GS.Name != "" {
		return nil
	}
	if plan.Bucket.Minio.Name != "" && plan.Bucket.Minio.Host != "" {
		return nil
	}

	return errors.New("missing S3, Minio bucket name or region or GS bucket name")
}
