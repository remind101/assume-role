package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"gopkg.in/yaml.v2"
	"runtime"
)

var configFilePath = fmt.Sprintf("%s/.aws/roles", os.Getenv("HOME"))

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s <role> [<command> <args...>]\n", os.Args[0])
	flag.PrintDefaults()
}

func init() {
	flag.Usage = usage
}

func main() {
	var (
		duration = flag.Duration("duration", time.Hour, "The duration that the credentials will be valid for.")
	)

	flag.Parse()
	argv := flag.Args()
	if len(argv) < 1 {
		flag.Usage()
		os.Exit(1)
	}

	stscreds.DefaultDuration = *duration

	role := argv[0]
	args := argv[1:]

	// Load credentials from configFilePath if it exists, else use regular AWS config
	var creds *credentials.Value
	if _, err := os.Stat(configFilePath); err == nil {
		fmt.Fprintf(os.Stderr, "WARNING: using deprecated role file (%s), switch to config file"+
			" (https://docs.aws.amazon.com/cli/latest/userguide/cli-roles.html)\n",
			configFilePath)
		config, err := loadConfig()
		must(err)

		roleConfig, ok := config[role]
		if !ok {
			must(fmt.Errorf("%s not in ~/.aws/roles", role))
		}

		// Clear out any previously set AWS_ environment variables so
		// they aren't used by this call
		cleanEnv()

		creds, err = assumeRole(roleConfig.Role, roleConfig.MFA, *duration)
		must(err)
	} else {
		cleanEnv()
		creds, err = assumeProfile(role)
		must(err)
	}

	if len(args) == 0 {
		if runtime.GOOS == "windows" {
			printWindowsCredentials(role, creds)
		} else {
			printCredentials(role, creds)
		}
		return
	}

	err := execWithCredentials(args, creds)
	must(err)
}

func cleanEnv() {
	os.Unsetenv("AWS_ACCESS_KEY_ID")
	os.Unsetenv("AWS_SECRET_ACCESS_KEY")
	os.Unsetenv("AWS_SESSION_TOKEN")
	os.Unsetenv("AWS_SECURITY_TOKEN")
}

func execWithCredentials(argv []string, creds *credentials.Value) error {
	argv0, err := exec.LookPath(argv[0])
	if err != nil {
		return err
	}

	os.Setenv("AWS_ACCESS_KEY_ID", creds.AccessKeyID)
	os.Setenv("AWS_SECRET_ACCESS_KEY", creds.SecretAccessKey)
	os.Setenv("AWS_SESSION_TOKEN", creds.SessionToken)
	os.Setenv("AWS_SECURITY_TOKEN", creds.SessionToken)

	env := os.Environ()
	return syscall.Exec(argv0, argv, env)
}

// printCredentials prints the credentials in a way that can easily be sourced
// with bash.
func printCredentials(role string, creds *credentials.Value) {
	fmt.Printf("export AWS_ACCESS_KEY_ID=\"%s\"\n", creds.AccessKeyID)
	fmt.Printf("export AWS_SECRET_ACCESS_KEY=\"%s\"\n", creds.SecretAccessKey)
	fmt.Printf("export AWS_SESSION_TOKEN=\"%s\"\n", creds.SessionToken)
	fmt.Printf("export AWS_SECURITY_TOKEN=\"%s\"\n", creds.SessionToken)
	fmt.Printf("export ASSUMED_ROLE=\"%s\"\n", role)
	fmt.Printf("# Run this to configure your shell:\n")
	fmt.Printf("# eval $(%s)\n", strings.Join(os.Args, " "))
}

// printWindowsCredentials prints the credentials in a way that can easily be sourced
// with Windows powershell using Invoke-Expression.
func printWindowsCredentials(role string, creds *credentials.Value) {
	fmt.Printf("$env:AWS_ACCESS_KEY_ID=\"%s\"\n", creds.AccessKeyID)
	fmt.Printf("$env:AWS_SECRET_ACCESS_KEY=\"%s\"\n", creds.SecretAccessKey)
	fmt.Printf("$env:AWS_SESSION_TOKEN=\"%s\"\n", creds.SessionToken)
	fmt.Printf("$env:AWS_SECURITY_TOKEN=\"%s\"\n", creds.SessionToken)
	fmt.Printf("$env:ASSUMED_ROLE=\"%s\"\n", role)
	fmt.Printf("# Run this to configure your shell:\n")
	fmt.Printf("# %s | Invoke-Expression \n", strings.Join(os.Args, " "))
}

// assumeProfile assumes the named profile which must exist in ~/.aws/config
// (https://docs.aws.amazon.com/cli/latest/userguide/cli-roles.html) and returns the temporary STS
// credentials.
func assumeProfile(profile string) (*credentials.Value, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Profile:                 profile,
		SharedConfigState:       session.SharedConfigEnable,
		AssumeRoleTokenProvider: readTokenCode,
	}))

	creds, err := sess.Config.Credentials.Get()
	if err != nil {
		return nil, err
	}
	return &creds, nil
}

// assumeRole assumes the given role and returns the temporary STS credentials.
func assumeRole(role, mfa string, duration time.Duration) (*credentials.Value, error) {
	sess := session.Must(session.NewSession())

	svc := sts.New(sess)

	params := &sts.AssumeRoleInput{
		RoleArn:         aws.String(role),
		RoleSessionName: aws.String("cli"),
		DurationSeconds: aws.Int64(int64(duration / time.Second)),
	}
	if mfa != "" {
		params.SerialNumber = aws.String(mfa)
		token, err := readTokenCode()
		if err != nil {
			return nil, err
		}
		params.TokenCode = aws.String(token)
	}

	resp, err := svc.AssumeRole(params)

	if err != nil {
		return nil, err
	}

	var creds credentials.Value
	creds.AccessKeyID = *resp.Credentials.AccessKeyId
	creds.SecretAccessKey = *resp.Credentials.SecretAccessKey
	creds.SessionToken = *resp.Credentials.SessionToken

	return &creds, nil
}

type roleConfig struct {
	Role string `yaml:"role"`
	MFA  string `yaml:"mfa"`
}

type config map[string]roleConfig

// readTokenCode reads the MFA token from Stdin.
func readTokenCode() (string, error) {
	r := bufio.NewReader(os.Stdin)
	fmt.Fprintf(os.Stderr, "MFA code: ")
	text, err := r.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(text), nil
}

// loadConfig loads the ~/.aws/roles file.
func loadConfig() (config, error) {
	raw, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}

	config := make(config)
	return config, yaml.Unmarshal(raw, &config)
}

func must(err error) {
	if err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			// Errors are already on Stderr.
			os.Exit(1)
		}

		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
