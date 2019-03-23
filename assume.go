package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"

	log "github.com/sirupsen/logrus"
)

var (
	duration  *time.Duration
	roleArnRe = regexp.MustCompile(`^arn:aws:iam::(.+):role/([^/]+)(/.+)?$`)
)

func init() {
	duration = flag.Duration("duration", time.Hour, "The duration that the credentials will be valid for.")
}

func loadCredentials(role string) (*credentials.Value, error) {
	// Load credentials from configFilePath if it exists, else use regular AWS config
	if roleArnRe.MatchString(role) {
		return assumeRole(role, "", *duration)
	}

	if _, err := os.Stat(configFilePath); err == nil {
		log.WithField("configFilePath", configFilePath).Warn(
			"Using deprecated role file, switch to config file" +
				" (https://docs.aws.amazon.com/cli/latest/userguide/cli-roles.html)")

		config, err := loadConfig()
		if err != nil {
			return nil, err
		}

		if roleConfig, ok := config[role]; ok {
			return assumeRole(roleConfig.Role, roleConfig.MFA, *duration)
		}

		return nil, fmt.Errorf("%s not in %s", role, configFilePath)
	}

	return assumeProfile(role)
}

// assumeProfile assumes the named profile which must exist in ~/.aws/config
// (https://docs.aws.amazon.com/cli/latest/userguide/cli-roles.html) and returns the temporary STS
// credentials.
func assumeProfile(profile string) (*credentials.Value, error) {
	log.Debug("Assuming role via named profile")
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
	log.Debug("Assume role via temporary STS credentials")
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
