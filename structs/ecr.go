package structs

// struct to hold ecr module data from terraform module
// this data if stored in the struct to allow the data to be queried for self approval
// data will be pulled in for a pull request and checked against the approvedModules map
type ECR struct {
	Source       string   `hcl:"source"`
	Name         string   `hcl:"repo_name"`
	OIDC         []string `hcl:"oidc_provider"`
	Github_Repos []string `hcl:"github_repositories"`
	Tags         Tags     `hcl:"tags"`
}
