package awsudo

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

// Options for awsudo.
type Options struct {
	// Role is the ARN of the role to assume.
	Role string
	// Profile from the AWS shared configuration to use. If unset, default is
	// used.
	Profile string
	// SessionName will be used as the Role Session Name during role assumption.
	SessionName string
}

// errEmptyCmd is returned the command passed to Command is empty.
var errEmptyCmd = errors.New("command may not be emtpy")

// Command prepares a command for execution in an environment with the
// credentials for an assumed AWS Role.
func Command(ctx context.Context, opts Options, command []string) (*exec.Cmd, error) {
	if len(command) < 1 {
		return nil, errEmptyCmd
	}

	creds, err := credentials(ctx, &opts)
	if err != nil {
		return nil, err
	}

	cmd := exec.CommandContext(ctx, command[0], command[1:]...)
	cmd.Env = prepareEnv(os.Environ(), &creds)

	return cmd, nil
}

// credentials uses the default AWS config for the environment to attempt to
// assume a role and retrieve credentials for it.
func credentials(ctx context.Context, opts *Options) (aws.Credentials, error) {
	var loadOptFns []func(*config.LoadOptions) error
	if opts.Profile != "" {
		loadOptFns = append(loadOptFns, config.WithSharedConfigProfile(opts.Profile))
	}

	awsCfg, err := config.LoadDefaultConfig(ctx, loadOptFns...)
	if err != nil {
		return aws.Credentials{}, err
	}

	credCache := aws.NewCredentialsCache(
		stscreds.NewAssumeRoleProvider(
			sts.NewFromConfig(awsCfg), opts.Role, withSessionName(opts.SessionName)))

	return credCache.Retrieve(ctx)
}

// withSessionName is an stscreds optFn which sets the RoleSessionName.
func withSessionName(session string) func(*stscreds.AssumeRoleOptions) {
	return func(opts *stscreds.AssumeRoleOptions) {
		opts.RoleSessionName = session
	}
}

// prepareEnv by pruning any existing AWS_* env vars, and setting the
// appropriate vars for the assumed role. Ensuring the env vars are set prevents
// the AWS SDK from attempting to look for creds elsewhere:
// https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-authentication.html#cli-chap-authentication-precedence
func prepareEnv(env []string, creds *aws.Credentials) []string {
	// Remove env vars from the ambient environment which might interfere with
	// AWS authentication in the executed command.
	env = filterPrefixes(env, []string{
		"AWS_ACCESS_KEY_ID",
		"AWS_PROFILE",
		"AWS_ROLE_ARN",
		"AWS_ROLE_SESSION_NAME",
		"AWS_SECRET_ACCESS_KEY",
		"AWS_SESSION_TOKEN",
	})

	// Now, repopulate the necessary env vars.
	env = append(env,
		"AWS_ACCESS_KEY_ID="+creds.AccessKeyID,
		"AWS_SECRET_ACCESS_KEY="+creds.SecretAccessKey,
		"AWS_SESSION_TOKEN="+creds.SessionToken,
	)

	return env
}

// filterPrefixes produces a slice of strings with any strings matching one or
// more of the provided prefixes removed.
//
// Yes, this is very inefficient, but there's little point being more clever for
// a set of strings as small as most environments.
func filterPrefixes(ss, prefixes []string) (ret []string) {
outerLoop:
	for _, s := range ss {
		for _, p := range prefixes {
			if strings.HasPrefix(s, p) {
				continue outerLoop
			}
		}
		ret = append(ret, s)
	}
	return
}
