package main

import (
	"log"

	"github.com/akhilsharma90/terraform-assistant/cmd/cli"
	"github.com/akhilsharma90/terraform-assistant/pkg/utils"
)

// It initializes the working directory and the Terraform executable directory,
// then calls the InitAndExecute function to start the program.
func main() {
	// Get the current working directory
	workingDir, err := utils.CurrenDir()
	if err != nil {
		log.Fatalf("Failed get current dir: %s\n", err.Error())
	}

	// Get the Terraform executable directory
	execDir, err := utils.TerraformPath()
	if err != nil {
		log.Fatalf("Failed get exec dir: %s\n", err.Error())
	}

	//the current directory and terraform executable directory that we have retrieved
	//by calling the helper functions above, we will pass to initandExecute function
	// Initialize and execute the program
	cli.InitAndExecute(workingDir, execDir)
}
