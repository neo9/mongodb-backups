package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

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

	return plan, nil
}

