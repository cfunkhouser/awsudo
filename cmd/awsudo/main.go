package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/cfunkhouser/awsudo"
	"github.com/spf13/pflag"
)

var (
	role = pflag.StringP("role", "u", "",
		"ARN for the AWS Role to assume for execution.")
	profile = pflag.StringP("profile", "p", "",
		"The AWS configuration profile to load. Empty for default.")
	session = pflag.StringP("session", "S", "awsudo",
		"The role session name to be used during role assumption.")
)

func init() {
	pflag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s - execute commands with AWS assumed role credentials\nUsage:\n", path.Base(os.Args[0]))
		pflag.PrintDefaults()
	}
}

func main() {
	pflag.Parse()

	args := pflag.Args()
	if len(args) < 1 {
		pflag.Usage()
		os.Exit(1)
	}

	if *role == "" {
		pflag.Usage()
		log.Fatal("role is required")
	}

	opts := awsudo.Options{
		Role:        *role,
		Profile:     *profile,
		SessionName: *session,
	}

	cmd, err := awsudo.Command(context.Background(), opts, args)
	if err != nil {
		log.Fatalf("Failed to prepare command execution: %v", err)
	}

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		// If the exec'ed process did not exit cleanly, pass along the code.
		os.Exit(cmd.ProcessState.ExitCode())
	}
}
