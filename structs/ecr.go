package structs

// ECR represents the configuration for an Elastic Container Registry (ECR).
// It includes the source of the ECR, the repository name, OIDC providers,
// associated GitHub repositories, and tags for the ECR.
//
// Fields:
// - Source: The source of the ECR.
// - Name: The name of the repository.
// - OIDC: A list of OIDC providers associated with the ECR.
// - Github_Repos: A list of GitHub repositories associated with the ECR.
// - Tags: Tags associated with the ECR.
type ECR struct {
	Source       string   `hcl:"source"`
	Name         string   `hcl:"repo_name"`
	OIDC         []string `hcl:"oidc_provider"`
	Github_Repos []string `hcl:"github_repositories"`
	Tags         Tags     `hcl:"tags"`
}
