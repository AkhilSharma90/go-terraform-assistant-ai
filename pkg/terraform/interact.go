package terraform

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

//FUNCTION NOT GETTING USED
// GetApplyConfirmation prompts the user for confirmation to apply changes.
// If requireConfirmation is false, it returns true without prompting the user.
// Otherwise, it displays a prompt asking the user to apply or not apply the changes.
// It returns true if the user selects "Apply", false if the user selects "Don't Apply",
// and an error if there is any issue with the prompt.
func GetApplyConfirmation(requireConfirmation bool) (bool, error) {
	if !requireConfirmation {
		return true, nil
	}

	prompt := promptui.Select{
		Label: "Would you like to apply this? [Apply/Don't Apply]",
		Items: []string{"Apply", "Don't Apply"},
	}

	_, result, err := prompt.Run()
	if err != nil {
		return false, fmt.Errorf("error prompt run: %w", err)
	}

	return result == "Apply", nil
}
