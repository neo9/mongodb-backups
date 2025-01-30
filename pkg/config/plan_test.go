//go:build !integration
// +build !integration

package config

import "testing"

func TestPlanOK(t *testing.T) {
	plan := Plan{}
	plan.GetPlan("plan.yaml")

	if plan.CreateDump.MaxRetries != DEFAULT_TRIES {
		t.Errorf("Expected %d obtained %d", DEFAULT_TRIES, plan.CreateDump.MaxRetries)
	}

	if plan.CreateDump.RetryDelay != DEFAULT_DELAY {
		t.Errorf("Expected %d obtained %d", DEFAULT_DELAY, plan.CreateDump.RetryDelay)
	}
}

func TestPlanKOFilename(t *testing.T) {
	plan := Plan{}
	_, err := plan.GetPlan("non-existing.yaml")
	if err == nil {
		t.Errorf("Expected error on the plan path")
	}
}

func TestPlanKOMultipleProviders(t *testing.T) {
	plan := Plan{}
	_, err := plan.GetPlan("plan_ko.yaml")
	if err == nil {
		t.Errorf("Expected error on file with minio and s3")
	}
}

func TestPlanKOJson(t *testing.T) {
	plan := Plan{}
	_, err := plan.GetPlan("plan_ko_format.yaml")
	if err == nil {
		t.Errorf("Expected error on the plan JSON format")
	}
}
