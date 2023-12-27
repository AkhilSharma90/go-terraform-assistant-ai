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

const (
	nameSubCommand = "You are a file name generator, only generate valid name for Terraform templates."
	runSubCommand  = "You are a Terraform HCL generator, only generate valid Terraform HCL without provider templates."
)

// runCommand is a function that executes the run command.
func runCommand(_ *cobra.Command, args []string) error {
	if len(args) == 0 {
		// If the length is 0, return an error with a wrapped error message
		return errors.Wrap(errLength, "prompt must be provided")
	}

	// Call the run function with the args
	return run(args)
}

// run is a function that executes the main logic of the CLI command.
//main business logic, takes care of everything from creating newOAIClients function
//to calling completion function to calling userPrompt function and many other helper functions
// It takes a slice of strings as input arguments and returns an error if any.
func run(args []string) error {
	// Create a context with a cancellation function that will be triggered on receiving an interrupt signal.
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// Create new OAI clients.
	oaiClients, err := newOAIClients()
	if err != nil {
		return fmt.Errorf("error creating new OAI client: %w", err)
	}

	var action, com, name string
	for action != apply {
		// Append the current action to the args slice.
		args = append(args, action)

		// Get completion for the run subcommand.
		//this creates the content for the terraform file
		com, err = completion(ctx, oaiClients, args, *openAIDeploymentName, runSubCommand)
		if err != nil {
			return fmt.Errorf("error completing run command: %w", err)
		}

		// Get completion for the name subcommand.
		//this just creates names of terraform files
		name, err = completion(ctx, oaiClients, args, *openAIDeploymentName, nameSubCommand)
		if err != nil {
			return fmt.Errorf("error completing name command: %w", err)
		}

		// Print the template to be stored.
		text := fmt.Sprintf("\n️🦄 Attempting to store the following template: %s", com)
		log.Println(text)

		// Prompt the user for an action.
		action, err = userActionPrompt()
		if err != nil {
			return err
		}

		// If the user chooses not to apply, return nil.
		if action == dontApply {
			return nil
		}
	}

	// Check the template for errors.
	if err = terraform.CheckTemplate(com); err != nil {
		return fmt.Errorf("error checking template: %w", err)
	}

	// Get the name from the completion result.
	name = utils.GetName(name)

	// Store the file with the given name and template.
	err = utils.StoreFile(name, com)
	if err != nil {
		return fmt.Errorf("error storing file: %w", err)
	}

	// Apply the Terraform operations.
	err = ops.Apply()
	if err != nil {
		return fmt.Errorf("error applying Terraform: %w", err)
	}

	return nil
}
