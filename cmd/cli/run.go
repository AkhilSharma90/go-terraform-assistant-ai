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

func runCommand(_ *cobra.Command, args []string) error {
	if len(args) == 0 {
		return errors.Wrap(errLength, "prompt must be provided")
	}

	return run(args)
}

func run(args []string) error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	oaiClients, err := newOAIClients()
	if err != nil {
		return fmt.Errorf("error new OAI client: %w", err)
	}

	var action, com, name string
	for action != apply {
		args = append(args, action)
		com, err = completion(ctx, oaiClients, args, *openAIDeploymentName, runSubCommand)

		if err != nil {
			return fmt.Errorf("error completion on run command: %w", err)
		}

		name, err = completion(ctx, oaiClients, args, *openAIDeploymentName, nameSubCommand)
		if err != nil {
			return fmt.Errorf("error completion on finding name command: %w", err)
		}

		text := fmt.Sprintf("\n️🦄 Attempting to store the following template: %s", com)
		log.Println(text)

		action, err = userActionPrompt()
		if err != nil {
			return err
		}

		if action == dontApply {
			return nil
		}
	}

	if err = terraform.CheckTemplate(com); err != nil {
		return fmt.Errorf("error check template: %w", err)
	}

	name = utils.GetName(name)

	err = utils.StoreFile(name, com)
	if err != nil {
		return fmt.Errorf("error store file: %w", err)
	}

	err = ops.Apply()
	if err != nil {
		return fmt.Errorf("error on apply terraform: %w", err)
	}

	return nil
}
