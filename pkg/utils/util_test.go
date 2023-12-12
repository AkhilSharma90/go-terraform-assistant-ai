package utils_test

import (
	"testing"

	"github.com/akhilsharma90/terraform-assistant/pkg/utils"
)

// TestRemoveBlankLinesFromString tests the RemoveBlankLinesFromString function.
// It verifies that the function correctly removes blank lines from a given string.
func TestRemoveBlankLinesFromString(t *testing.T) {
	// Define the input string with multiple blank lines.
	input := "\n\n\nHello, world!\n\nHow are you?\n\n\n"

	// Define the expected output string with blank lines removed.
	expectedOutput := "Hello, world!\n\nHow are you?\n\n\n"

	// Call the RemoveBlankLinesFromString function to get the actual output.
	output := utils.RemoveBlankLinesFromString(input)

	// Compare the actual output with the expected output.
	if output != expectedOutput {
		t.Errorf("Expected output '%s', but got '%s'", expectedOutput, output)
	}
}
