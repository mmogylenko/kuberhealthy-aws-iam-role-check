package main

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Setenv("TARGET_ARN", "arn:aws:sts::012345678910:role/read-access-to-everything")
	targetArnEnv = "arn:aws:sts::012345678910:role/read-access-to-everything"
	os.Exit(m.Run())
}
