package main

import (
	"flag"
	"os"
)

var reset *bool

func init() {
	reset = flag.Bool("reset", false, "Reset AWS Env-var tokens (internally) before retrieving new credentials.")
}

func resetEnvVars() {
	os.Unsetenv("AWS_ACCESS_KEY_ID")
	os.Unsetenv("AWS_SECRET_ACCESS_KEY")
	os.Unsetenv("AWS_SESSION_TOKEN")
	os.Unsetenv("AWS_SECURITY_TOKEN")
	os.Unsetenv("ASSUMED_ROLE")
}
