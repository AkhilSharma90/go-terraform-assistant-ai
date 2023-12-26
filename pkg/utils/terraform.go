package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)
//NOT GETTING USED
func EndsWithTf(str string) bool {
	return strings.HasSuffix(str, ".tf")
}
//NOT GETTING USED
func RandomName() string {
	// Initialize a byte slice of desired length
	randomBytes := make([]byte, 5)

	// Read random data into the byte slice
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}

	// Encode the byte slice as a base64 string
	randomString := base64.RawURLEncoding.EncodeToString(randomBytes)

	// Truncate the string to the desired length

	return fmt.Sprintf("terraform-%s.tf", randomString)
}

//Getting called in the run function in run.go
// GetName returns a modified version of the input name string.
func GetName(name string) string {
	name = RemoveBlankLinesFromString(name)
	if EndsWithTf(name) {
		return name
	}

	return RandomName()
}

//Getting called from the main function
// TerraformPath returns the path of the Terraform executable.
// It uses the "where" command on Windows and the "which" command on other platforms to locate the Terraform executable.
// Returns the path of the Terraform executable and any error encountered.
func TerraformPath() (string, error) {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		// Use the "where" command to locate the Terraform executable on Windows.
		cmd = exec.Command("where", "terraform")
	} else {
		// Use the "which" command to locate the Terraform executable on other platforms.
		cmd = exec.Command("which", "terraform")
	}

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("error running Init: %w", err)
	}

	// Trim the newline character from the output and return the path of the Terraform executable.
	return strings.TrimRight(string(output), "\n"), nil
}
