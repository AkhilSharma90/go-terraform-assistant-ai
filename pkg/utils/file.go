package utils

import (
	"fmt"
	"log"
	"os"
)

// DirExists checks if a directory exists at the specified path.
// It returns true if the directory exists, and false otherwise.
//FUNCTION NOT USED
func DirExists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		log.Println("Terraform is not initialized. Run `terraform init` first.")
	} else if err != nil {
		log.Fatalf("Failed to check if Terraform is initialized: %s\n", err.Error())
	}

	return true
}

//Getting called from initCommand function
// StoreFile writes the contents to a file with the given name.
// It removes blank lines from the contents before writing.
func StoreFile(name string, contents string) error {
	contents = RemoveBlankLinesFromString(contents)

	err := os.WriteFile(name, []byte(contents), 0o600)
	if err != nil {
		return fmt.Errorf("error writing file: %w", err)
	}

	return nil
}

//Getting called from the main function
// CurrenDir returns the current working directory.
func CurrenDir() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("error current dir: %w", err)
	}

	return currentDir, nil
}
