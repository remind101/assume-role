package main

import (
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"
)

func must(err error) {
	if err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			// Errors are already on Stderr.
			os.Exit(1)
		}

		log.Errorf("%v", err)
		os.Exit(1)
	}
}
