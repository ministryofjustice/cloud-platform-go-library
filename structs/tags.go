package structs

// Tags is a struct that contains the standard tags used across all Cloud Plaform Modules.
// Each field corresponds to a specific tag with its respective HCL (HashiCorp Configuration Language) key.
// Feilds:
// - Business_unit: The business unit associated with the resource.
// - Application: The application name associated with the resource.
// - Is_production: Indicates whether the resource is part of a production environment.
// - Team_name: The name of the team responsible for the resource.
// - Namespace: The namespace in which the resource resides.
// - Environment_name: The name of the environment (e.g., dev, staging, prod) associated with the resource.
// - Infrastructure_support: The contact information or team responsible for infrastructure support.
type Tags struct {
	Business_unit          string `hcl:"business_unit"`
	Application            string `hcl:"application"`
	Is_production          string `hcl:"is_production"`
	Team_name              string `hcl:"team_name"`
	Namespace              string `hcl:"namespace"`
	Environment_name       string `hcl:"environment_name"`
	Infrastructure_support string `hcl:"infrastructure_support"`
}
