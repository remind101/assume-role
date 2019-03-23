package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	log "github.com/sirupsen/logrus"
)

var (
	debug *bool
	// verbose *bool
)

func init() {
	debug = flag.Bool("debug", false, "Show debug logs")
	// verbose = flag.Bool("verbose", false, "Increased verbosity")
	flag.Usage = usage
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [<args...>] <role> [<command> <command args...>]\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	flag.Parse()
	if *debug {
		log.SetLevel(log.DebugLevel)
	}

	log.Debug("Starting...")
	argv := flag.Args()
	if len(argv) < 1 {
		flag.Usage()
		os.Exit(1)
	}

	log.WithField("duration", *duration).Debug("Setting STS Token Duration")
	stscreds.DefaultDuration = *duration

	role := argv[0]
	args := argv[1:]

	if *reset {
		resetEnvVars()
	}

	creds, err := loadCredentials(role)
	must(err)
	log.WithField("credentials", creds).Debug("Received role credentials")

	if len(args) == 0 {
		printCredentials(role, creds)
		return
	}

	must(execWithCredentials(role, args, creds))
}
