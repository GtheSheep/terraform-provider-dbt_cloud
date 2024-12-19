package resources_test

import (
	"fmt"
	"github.com/dbt-labs/terraform-provider-dbtcloud/pkg/framework/acctest_config"
	"regexp"
	"strings"
	"testing"

	"github.com/dbt-labs/terraform-provider-dbtcloud/pkg/framework/acctest_helper"
	"github.com/dbt-labs/terraform-provider-dbtcloud/pkg/helper"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TODO: Add more tests
// we are currently testing in CI the SSH cloning and GH native cloning but not GitLab native and ADO native
// this would require having the GitLab and ADO native integrations set up in the dbt Cloud account used for CI

func TestAccDbtCloudRepositoryResource(t *testing.T) {

	repoUrlGithub := "git@github.com:dbt-labs/terraform-provider-dbtcloud.git"
	projectName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest_helper.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDbtCloudRepositoryDestroy,
		Steps: []resource.TestStep{
			// Create Github repository
			{
				Config: testAccDbtCloudRepositoryResourceGithubConfig(repoUrlGithub, projectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDbtCloudRepositoryExists(
						"dbtcloud_repository.test_repository_github",
					),
					resource.TestCheckResourceAttr(
						"dbtcloud_repository.test_repository_github",
						"remote_url",
						repoUrlGithub,
					),
					resource.TestCheckResourceAttrSet(
						"dbtcloud_repository.test_repository_github",
						"deploy_key",
					),
				),
			},
			// MODIFY
			// IMPORT
			{
				ResourceName:            "dbtcloud_repository.test_repository_github",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"fetch_deploy_key"},
			},
		},
	})

	repoUrlGithubApplication := acctest_config.AcceptanceTestConfig.GitHubRepoUrl
	githubAppInstallationId := acctest_config.AcceptanceTestConfig.GitHubAppInstallationId
	projectNameGithubApplication := strings.ToUpper(
		acctest.RandStringFromCharSet(10, acctest.CharSetAlpha),
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest_helper.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDbtCloudRepositoryDestroy,
		Steps: []resource.TestStep{
			// Create Github repository via the GithUb Application
			{
				Config: testAccDbtCloudRepositoryResourceGithubApplicationConfig(
					repoUrlGithubApplication,
					projectNameGithubApplication,
					githubAppInstallationId,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDbtCloudRepositoryExists(
						"dbtcloud_repository.test_repository_github_application",
					),
					resource.TestCheckResourceAttr(
						"dbtcloud_repository.test_repository_github_application",
						"remote_url",
						repoUrlGithubApplication,
					),
					resource.TestCheckResourceAttr(
						"dbtcloud_repository.test_repository_github_application",
						"git_clone_strategy",
						"github_app",
					),
				),
			},
			// MODIFY
			// IMPORT
			{
				ResourceName:            "dbtcloud_repository.test_repository_github_application",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"fetch_deploy_key"},
			},
		},
	})
}

func testAccDbtCloudRepositoryResourceGithubConfig(repoUrl, projectName string) string {
	return fmt.Sprintf(`
resource "dbtcloud_project" "test_project" {
  name        = "%s"
}

resource "dbtcloud_repository" "test_repository_github" {
  remote_url = "%s"
  project_id = dbtcloud_project.test_project.id
  git_clone_strategy = "deploy_key"
  pull_request_url_template = "https://github.com/my-org/my-repo/compare/qa...{{source}}"
}
`, projectName, repoUrl)
}

func testAccDbtCloudRepositoryResourceGithubApplicationConfig(
	repoUrl string,
	projectName string,
	githubAppInstallationId int,
) string {
	return fmt.Sprintf(`
resource "dbtcloud_project" "test_project" {
  name        = "%s"
}

resource "dbtcloud_repository" "test_repository_github_application" {
  remote_url = "%s"
  project_id = dbtcloud_project.test_project.id
  github_installation_id = %d
  git_clone_strategy = "github_app"
  pull_request_url_template = "https://github.com/my-org/my-repo/compare/qa...{{source}}"
}
`, projectName, repoUrl, githubAppInstallationId)
}

//
// func testAccDbtCloudRepositoryResourceGitlabConfig(repoUrl, projectName string) string {
// 	return fmt.Sprintf(`
// resource "dbtcloud_project" "test_project" {
//   name        = "%s"
// }
//
// resource "dbtcloud_repository" "test_repository_gitlab" {
//   remote_url = "%s"
//   project_id = dbtcloud_project.test_project.id
//   gitlab_project_id = 34786716
// }
// `, projectName, repoUrl)
// }

func testAccCheckDbtCloudRepositoryExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}
		apiClient, err := acctest_helper.SharedClient()
		if err != nil {
			return fmt.Errorf("Issue getting the client")
		}
		projectId, repositoryId, err := helper.SplitIDToStrings(
			rs.Primary.ID,
			"dbtcloud_repository",
		)
		if err != nil {
			return err
		}
		_, err = apiClient.GetRepository(repositoryId, projectId)
		if err != nil {
			return fmt.Errorf("error fetching item with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccCheckDbtCloudRepositoryDestroy(s *terraform.State) error {
	apiClient, err := acctest_helper.SharedClient()
	if err != nil {
		return fmt.Errorf("Issue getting the client")
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "dbtcloud_repository" {
			continue
		}
		projectId, repositoryId, err := helper.SplitIDToStrings(
			rs.Primary.ID,
			"dbtcloud_repository",
		)
		if err != nil {
			return err
		}
		_, err = apiClient.GetRepository(repositoryId, projectId)
		if err == nil {
			return fmt.Errorf("Repository still exists")
		}
		notFoundErr := "resource-not-found"
		expectedErr := regexp.MustCompile(notFoundErr)
		if !expectedErr.Match([]byte(err.Error())) {
			return fmt.Errorf("expected %s, got %s", notFoundErr, err)
		}
	}

	return nil
}
