package terraform

import (
	"context"
	"fmt"
	"time"

	"github.com/briandowns/spinner"
)

// Init initializes the Terraform instance.
// It starts a spinner, runs the Init command, and stops the spinner.
// Returns an error if there was an error running Init.
func (ter *Terraform) Init() error {
	spin := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	spin.Start()

	err := ter.Exec.Init(context.Background())
	if err != nil {
		spin.Stop()

		return fmt.Errorf("error running Init: %w", err)
	}

	spin.Stop()

	return nil
}

// Apply applies the Terraform configuration.
// It starts a spinner to indicate that the apply process is running.
// It then calls the Exec.Apply method to execute the apply command.
func (ter *Terraform) Apply() error {
	spin := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	spin.Start()

	err := ter.Exec.Apply(context.Background())
	if err != nil {
		spin.Stop()

		return fmt.Errorf("error running Apply: %w", err)
	}

	spin.Stop()

	return nil
}
