package structs

type ECR struct {
	Source       string   `hcl:"source"`
	Name         string   `hcl:"repo_name"`
	OIDC         []string `hcl:"oidc_provider"`
	Github_Repos []string `hcl:"github_repositories"`
	Tags         Tags     `hcl:"tags"`
}
