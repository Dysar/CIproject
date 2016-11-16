package ciproject_test

import (
	"testing"
	ciproject "."
)

func TestResult(t *testing.T) {
	if err := ciproject.Result(20, 15) < 50; err != true {
		t.Error("test cannot fail")
	}
}

func TestResultFail(t *testing.T) {
	if err := ciproject.Result(20, 15) > 30; err != true {
		t.Error("GG man")
	}
}