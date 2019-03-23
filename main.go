package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [<args...>] <role> [<command> <command args...>]\n", os.Args[0])
	flag.PrintDefaults()
}

func init() {
	flag.Usage = usage
}

func main() {
	flag.Parse()
	argv := flag.Args()
	if len(argv) < 1 {
		flag.Usage()
		os.Exit(1)
	}

	stscreds.DefaultDuration = *duration

	role := argv[0]
	args := argv[1:]

	if *reset {
		resetEnvVars()
	}

	creds, err := loadCredentials(role)
	must(err)

	if len(args) == 0 {
		printCredentials(role, creds)
		return
	}

	must(execWithCredentials(role, args, creds))
}
