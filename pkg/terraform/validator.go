package terraform

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/pkg/errors"
)

var errTemplate = errors.New("invalid terraform template")

//Getting called from initCmd function
// CheckTemplate validates the given completion string as a Terraform template.
// It checks if the template is valid by parsing it using hclsyntax.ParseConfig.
// If any parse diagnostics are found, an error is returned.
func CheckTemplate(completion string) error {
	template := []byte(completion)
	_, parseDiags := hclsyntax.ParseConfig(template, "", hcl.Pos{Line: 2, Column: 1})

	if len(parseDiags) != 0 {
		return errors.Wrapf(errTemplate, "expected valid template but: %s", completion)
	}

	return nil
}
