package structs

type Tags struct {
	Business_unit          string `hcl:"business_unit"`
	Application            string `hcl:"application"`
	Is_production          string `hcl:"is_production"`
	Team_name              string `hcl:"team_name"`
	Namespace              string `hcl:"namespace"`
	Environment_name       string `hcl:"environment_name"`
	Infrastructure_support string `hcl:"infrastructure_support"`
}
