package terraform

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	gohcl "github.com/hashicorp/hcl/v2/gohcl"
	hclparse "github.com/hashicorp/hcl/v2/hclparse"
	structs "github.com/ministryofjustice/cloud-platform-go-library/structs"
)

// GetTFBody parses an HCL file from the given source path and returns its body.
//
// Parameters:
//   - source: A string representing the file path to the HCL file.
//
// Returns:
//   - hcl.Body: The body of the parsed HCL file.
//   - error: An error if the file could not be parsed, otherwise nil.
func GetTFBody(source string) (hcl.Body, error) {
	parser := hclparse.NewParser()
	file, diags := parser.ParseHCLFile(source)
	if diags.HasErrors() {
		return nil, fmt.Errorf("error parsing HCL file")
	}
	return file.Body, nil
}

// MapTfFileToStruct maps a Terraform file to a corresponding Go struct based on the provided source.
// It decodes the HCL body into the appropriate struct.
//
// Parameters:
//   - source: A string representing the source module.
//   - body: An hcl.Body containing the HCL content to be decoded.
//
// Returns:
//   - An interface{} containing the decoded struct if successful.
//   - An error if the decoding fails or if the source module is not found.
func MapTfFileToStruct(source string, body hcl.Body) (interface{}, error) {
	switch source {
	case "github.com/ministryofjustice/cloud-platform-terraform-ecr-credentials":
		var ecr structs.ECR
		diags := gohcl.DecodeBody(body, nil, &ecr)
		if diags.HasErrors() {
			fmt.Println("error decoding HCL file")
			return nil, fmt.Errorf("error decoding HCL file")
		}
		return ecr, nil
	default:
		return nil, fmt.Errorf("module not found")
	}
}
