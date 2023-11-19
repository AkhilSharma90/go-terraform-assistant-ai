package utils_test

import (
	"testing"

	"github.com/akhilsharma90/terraform-assistant/pkg/utils"
)

func TestRemoveBlankLinesFromString(t *testing.T) {
	input := "\n\n\nHello, world!\n\nHow are you?\n\n\n"
	expectedOutput := "Hello, world!\n\nHow are you?\n\n\n"

	output := utils.RemoveBlankLinesFromString(input)

	if output != expectedOutput {
		t.Errorf("Expected output '%s', but got '%s'", expectedOutput, output)
	}
}
