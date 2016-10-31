package main_test

import (
	"testing"
	ciproject "../CIproject"
)

func TestResult(t *testing.T) {
	if ciproject.Result(20, 15) > 50 {
		t.Error("OP man")
	}
}

func TestResultFail(t *testing.T) {
	if ciproject.Result(20, 15) < 50 {
		t.Error("GG man")
	}
}