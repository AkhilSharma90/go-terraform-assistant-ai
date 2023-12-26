package terraform

import (
	"fmt"

	"github.com/hashicorp/terraform-exec/tfexec"
)

type Terraform struct {
	WorkingDir string
	ExecDir    string
	Exec       *tfexec.Terraform
}


//Getting called from rootCmd function
// NewTerraform creates a new instance of the Terraform struct.
func NewTerraform(workingDir string, execDir string) (*Terraform, error) {
	// Create a new instance of tfexec.Terraform using the provided working directory and execution directory.
	tf, err := tfexec.NewTerraform(workingDir, execDir)
	if err != nil {
		return nil, fmt.Errorf("error new terraform: %w", err)
	}

	// Create a new Terraform struct with the provided working directory, execution directory, and tfexec.Terraform instance.
	return &Terraform{
		WorkingDir: workingDir,
		ExecDir:    execDir,
		Exec:       tf,
	}, nil
}
