package main

import (
	"flag"
	"os"

	log "github.com/sirupsen/logrus"
)

var reset *bool

func init() {
	reset = flag.Bool("reset", false, "If a profile is provided: internally reset AWS Env-var tokens before retrieving new credentials.\nIf not: output the Env-var reset commands.")
}

func resetEnvVars() {
	log.WithFields(log.Fields{
		"AWS_ACCESS_KEY_ID":     os.Getenv("AWS_ACCESS_KEY_ID"),
		"AWS_SECRET_ACCESS_KEY": os.Getenv("AWS_SECRET_ACCESS_KEY"),
		"AWS_SESSION_TOKEN":     os.Getenv("AWS_SESSION_TOKEN"),
		"AWS_SECURITY_TOKEN":    os.Getenv("AWS_SECURITY_TOKEN"),
		"ASSUMED_ROLE":          os.Getenv("ASSUMED_ROLE"),
	}).Debug("Resetting Token Env-vars. Showing prev vals")

	os.Unsetenv("AWS_ACCESS_KEY_ID")
	os.Unsetenv("AWS_SECRET_ACCESS_KEY")
	os.Unsetenv("AWS_SESSION_TOKEN")
	os.Unsetenv("AWS_SECURITY_TOKEN")
	os.Unsetenv("ASSUMED_ROLE")
}
