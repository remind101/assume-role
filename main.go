package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"gopkg.in/yaml.v2"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws"
)

var configFilePath = fmt.Sprintf("%s/.aws/roles", os.Getenv("HOME"))

func usage() {
	fmt.Print(`Usage: assume-role <role> [<command> <args...>]
`)
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	role := os.Args[1]
	args := os.Args[2:]

	config, err := loadConfig()
	must(err)

	roleConfig, ok := config[role]
	if !ok {
		must(fmt.Errorf("%s not in ~/.aws/roles", role))
	}

	if os.Getenv("ASSUMED_ROLE") != "" {
		// Clear out any previously set AWS_ environment variables so
		// they aren't used with the assumeRole command.
		cleanEnv()
	}

	creds, err := assumeRole(roleConfig.Role, roleConfig.MFA)
	must(err)

	if len(args) == 0 {
		printCredentials(role, creds)
		return
	}

	err = execWithCredentials(args, creds)
	must(err)
}

func cleanEnv() {
	os.Unsetenv("AWS_ACCESS_KEY_ID")
	os.Unsetenv("AWS_SECRET_ACCESS_KEY")
	os.Unsetenv("AWS_SESSION_TOKEN")
	os.Unsetenv("AWS_SECURITY_TOKEN")
}

func execWithCredentials(argv []string, creds *credentials) error {
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

type credentials struct {
	AccessKeyID     string
	SecretAccessKey string
	SessionToken    string
}

// printCredentials prints the credentials in a way that can easily be sourced
// with bash.
func printCredentials(role string, creds *credentials) {
	fmt.Printf("export AWS_ACCESS_KEY_ID=\"%s\"\n", creds.AccessKeyID)
	fmt.Printf("export AWS_SECRET_ACCESS_KEY=\"%s\"\n", creds.SecretAccessKey)
	fmt.Printf("export AWS_SESSION_TOKEN=\"%s\"\n", creds.SessionToken)
	fmt.Printf("export AWS_SECURITY_TOKEN=\"%s\"\n", creds.SessionToken)
	fmt.Printf("export ASSUMED_ROLE=\"%s\"\n", role)
	fmt.Printf("# Run this to configure your shell:\n")
	fmt.Printf("# eval $(%s)\n", strings.Join(os.Args, " "))
}

// assumeRole assumes the given role and returns the temporary STS credentials.

func assumeRole(role, mfa string) (*credentials, error) {
	sess := session.Must(session.NewSession())

	svc := sts.New(sess)

	params := &sts.AssumeRoleInput{
		RoleArn: aws.String(role),
		RoleSessionName: aws.String("cli"),
	}
	if mfa != "" {
		params.SerialNumber = aws.String(mfa)
		params.TokenCode = aws.String(readTokenCode())
	}

	resp, err := svc.AssumeRole(params)

	if err != nil {
		return nil, err
	}

	var creds credentials
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
func readTokenCode() string {
	r := bufio.NewReader(os.Stdin)
	fmt.Fprintf(os.Stderr, "MFA code: ")
	text, _ := r.ReadString('\n')
	return strings.TrimSpace(text)
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
