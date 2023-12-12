package cli

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/akhilsharma90/terraform-assistant/pkg/terraform"
	"github.com/akhilsharma90/terraform-assistant/pkg/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// Constant string for the init subcommand description
const initSubCommand = "You are a Terraform HCL generator, only generate valid provider Terraform HCL templates."

// Error for invalid length
var errLength = errors.New("invalid length")

// addInit creates and returns a new Cobra command for the "init" subcommand.
// This command is used to run the "terraform init" command.
func addInit() *cobra.Command {
	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Run terraform init",
		RunE:  initCommand,
	}

	return initCmd
}

// initCommand is a function that handles the "init" command in the CLI.
// The function checks if the length of the args slice is 0 and returns an error if it is.
func initCommand(_ *cobra.Command, args []string) error {
	if len(args) == 0 {
		return errors.Wrap(errLength, "prompt must be provided")
	}

	return initCmd(args)
}

// initCmd initializes the command for initializing the Terraform project.
// It creates the necessary files, checks the template, and runs Terraform init.
func initCmd(args []string) error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// Create new OAI clients
	oaiClients, err := newOAIClients()
	if err != nil {
		return fmt.Errorf("error creating new OAI client: %w", err)
	}

	var action, com string
	for action != apply {
		args = append(args, action)

		// Get completion for the current command
		com, err = completion(ctx, oaiClients, args, *openAIDeploymentName, initSubCommand)
		if err != nil {
			return fmt.Errorf("error completion: %w", err)
		}

		text := fmt.Sprintf("\n🦄 Attempting to apply the following template: %s", com)
		log.Println(text)

		// Prompt user for action
		action, err = userActionPrompt()
		if err != nil {
			return err
		}

		if action == dontApply {
			return nil
		}
	}

	// Check the template
	if err = terraform.CheckTemplate(com); err != nil {
		return fmt.Errorf("error checking template: %w", err)
	}

	// Store the template in a file
	if err = utils.StoreFile("provider.tf", com); err != nil {
		return fmt.Errorf("error storing file: %w", err)
	}

	// Run Terraform init
	if err = ops.Init(); err != nil {
		return fmt.Errorf("error running terraform init: %w", err)
	}

	return nil
}
