package cli

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

const (
	apply     = "Apply"
	dontApply = "Don't Apply"
	reprompt  = "Reprompt"
)

// userActionPrompt prompts the user for an action and returns the selected action.
func userActionPrompt() (string, error) {
	var (
		result string
		err    error
	)

	// If requireConfirmation flag is not set, return the default action as apply
	if !*requireConfirmation {
		return apply, nil
	}

	// Create a label for the prompt
	items := []string{apply, dontApply}
	label := fmt.Sprintf("Would you like to apply this? [%s/%s/%s]", reprompt, items[0], items[1])

	// Create a prompt with options to select apply or not apply
	prompt := promptui.SelectWithAdd{
		Label:    label,
		Items:    items,
		AddLabel: reprompt,
	}

	// Run the prompt and get the selected action
	_, result, err = prompt.Run()

	if err != nil {
		return dontApply, fmt.Errorf("error to run prompt: %w", err)
	}

	return result, nil
}
