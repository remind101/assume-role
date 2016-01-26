package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"gopkg.in/yaml.v2"
)

var configFilePath = fmt.Sprintf("%s/.aws/roles", os.Getenv("HOME"))

func main() {
	role := os.Args[1]
	args := os.Args[2:]

	config, err := loadConfig()
	must(err)

	roleConfig, ok := config[role]
	if !ok {
		must(fmt.Errorf("%s not in ~/.aws/roles", role))
	}

	creds, err := assumeRole(roleConfig.Role, roleConfig.MFA)
	must(err)

	if len(args) == 0 {
		printCredentials(creds)
		return
	}

	err = execWithCredentials(args, creds)
	must(err)
}

func execWithCredentials(argv []string, creds *credentials) error {
	argv0, err := exec.LookPath(argv[0])
	if err != nil {
		return err
	}

	env := os.Environ()
	env = append(env, fmt.Sprintf("AWS_ACCESS_KEY_ID=%s", creds.AccessKeyId))
	env = append(env, fmt.Sprintf("AWS_SECRET_ACCESS_KEY=%s", creds.SecretAccessKey))
	env = append(env, fmt.Sprintf("AWS_SESSION_TOKEN=%s", creds.SessionToken))
	return syscall.Exec(argv0, argv, env)
}

type credentials struct {
	AccessKeyId     string
	SecretAccessKey string
	SessionToken    string
}

// printCredentials prints the credentials in a way that can easily be sourced
// with bash.
func printCredentials(creds *credentials) {
	fmt.Printf("export AWS_ACCESS_KEY_ID=\"%s\"\n", creds.AccessKeyId)
	fmt.Printf("export AWS_SECRET_ACCESS_KEY=\"%s\"\n", creds.SecretAccessKey)
	fmt.Printf("export AWS_SESSION_TOKEN=\"%s\"\n", creds.SessionToken)
	fmt.Printf("# Run this to configure your shell:\n")
	fmt.Printf("# eval $(%s)\n", strings.Join(os.Args, " "))
}

// assumeRole assumes the given role and returns the temporary STS credentials.
func assumeRole(role, mfa string) (*credentials, error) {
	args := []string{
		"sts",
		"assume-role",
		"--output", "json",
		"--role-arn", role,
		"--role-session-name", "cli",
	}
	if mfa != "" {
		args = append(args,
			"--serial-number", mfa,
			"--token-code",
			readTokenCode(),
		)
	}

	b := new(bytes.Buffer)
	cmd := exec.Command("aws", args...)
	cmd.Stdout = b
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	if err := cmd.Wait(); err != nil {
		return nil, err
	}

	var resp struct{ Credentials credentials }
	if err := json.NewDecoder(b).Decode(&resp); err != nil {
		return nil, err
	}

	return &resp.Credentials, nil
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
