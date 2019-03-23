package main

import (
	"os"
	"os/exec"
	"syscall"

	"github.com/aws/aws-sdk-go/aws/credentials"
)

func execWithCredentials(role string, argv []string, creds *credentials.Value) error {
	argv0, err := exec.LookPath(argv[0])
	if err != nil {
		return err
	}

	os.Setenv("AWS_ACCESS_KEY_ID", creds.AccessKeyID)
	os.Setenv("AWS_SECRET_ACCESS_KEY", creds.SecretAccessKey)
	os.Setenv("AWS_SESSION_TOKEN", creds.SessionToken)
	os.Setenv("AWS_SECURITY_TOKEN", creds.SessionToken)
	os.Setenv("ASSUMED_ROLE", role)

	env := os.Environ()
	return syscall.Exec(argv0, argv, env)
}
