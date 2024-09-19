package terraform

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	gohcl "github.com/hashicorp/hcl/v2/gohcl"
	hclparse "github.com/hashicorp/hcl/v2/hclparse"
	structs "github.com/ministryofjustice/cloud-platform-go-library/structs"
)

func GetTFBody(source string) (hcl.Body, error) {
	// Parse the HCL file
	parser := hclparse.NewParser()
	file, diags := parser.ParseHCLFile(source)
	if diags.HasErrors() {
		fmt.Println("error parsing HCL file")
		return nil, fmt.Errorf("error parsing HCL file")
	}
	return file.Body, nil
}

func MapTFFile(source string, body hcl.Body) (interface{}, error) {
	// switch case to map the terraform file to the struct
	// this will allow the data to be queried for self approval

	switch source {
	case "github.com/ministryofjustice/cloud-platform-terraform-ecr-credentials":
		var ecr structs.ECR
		// Decode the HCL file
		diags := gohcl.DecodeBody(body, nil, &ecr)
		if diags.HasErrors() {
			fmt.Println("error decoding HCL file")
			return nil, fmt.Errorf("error decoding HCL file")
		}
		return ecr, nil
	default:
		fmt.Println("module not found")
		return nil, fmt.Errorf("module not found")
	}
}
