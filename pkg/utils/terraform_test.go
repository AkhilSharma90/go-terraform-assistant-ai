package utils_test

import (
	"testing"

	"github.com/akhilsharma90/terraform-assistant/pkg/utils"
)

// TestEndsWithTf tests the EndsWithTf function to check if a given string ends with ".tf".
func TestEndsWithTf(t *testing.T) {
	cases := []struct {
		input    string
		expected bool
	}{
		{"main.tf", true},
		{"provider.txt", false},
		{"provider.tf", true},
		{"subnet.doc", false},
	}

	for _, c := range cases {
		result := utils.EndsWithTf(c.input)

		if result != c.expected {
			t.Errorf("EndsWithTf(%q) == %v, expected %v", c.input, result, c.expected)
		}
	}
}

// TestRandomName is a unit test function that tests the RandomName function in the utils package.
// It generates two random names and checks if they are unique.
func TestRandomName(t *testing.T) {
	name1 := utils.RandomName()
	name2 := utils.RandomName()

	if name1 == name2 {
		t.Errorf("Expected unique names, but got: %s, %s", name1, name2)
	}
}
