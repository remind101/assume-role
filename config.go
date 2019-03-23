package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

type roleConfig struct {
	Role string `yaml:"role"`
	MFA  string `yaml:"mfa"`
}

type config map[string]roleConfig

var configFilePath = fmt.Sprintf("%s/.aws/roles", os.Getenv("HOME"))

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

	roleConfig := make(config)
	return roleConfig, yaml.Unmarshal(raw, &roleConfig)
}
