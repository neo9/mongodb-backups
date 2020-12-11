package config

import (
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type Plan struct {
	Name string `json:"name"`
	Schedule string `json:"schedule"`
	Retention string `json:"retention"`
	Timeout string `json:"timeout"`
	MongoDB MongoDB `json:"mongodb"`
	Bucket Bucket `json:"buckets"`
}

type Bucket struct {
	S3 S3 `json:"s3"`
	GS GS `json:"gs"`
}

type S3 struct {
	Name string `json:"name"`
	Region string `json:"region"`
}
type GS struct {
	Name string `json:"name"`
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
	if plan.Bucket.S3.Name != "" && plan.Bucket.GS.Name != "" {
		return errors.New("error in configuration : should only have s3 OR gs bucket configured")
	}

	if plan.Bucket.S3.Name != "" && plan.Bucket.S3.Region != "" {
		return nil
	}
	if plan.Bucket.GS.Name != "" {
		return nil
	}

	return errors.New("missing S3 bucket name or region or GS bucket name")
}
