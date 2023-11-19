package main

import (
	"log"

	"github.com/akhilsharma90/terraform-assistant/cmd/cli"
	"github.com/akhilsharma90/terraform-assistant/pkg/utils"
)

func main() {
	workingDir, err := utils.CurrenDir()
	if err != nil {
		log.Fatalf("Failed get current dir: %s\n", err.Error())
	}

	execDir, err := utils.TerraformPath()
	if err != nil {
		log.Fatalf("Failed get exec dir: %s\n", err.Error())
	}

	cli.InitAndExecute(workingDir, execDir)
}
