package resources_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/gthesheep/terraform-provider-dbt-cloud/pkg/dbt_cloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDbtCloudProjectRepositoryResource(t *testing.T) {

	projectName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	repoUrlGithub := "git@github.com:GtheSheep/terraform-provider-dbt-cloud.git"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDbtCloudProjectRepositoryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDbtCloudProjectRepositoryResourceBasicConfig(projectName, repoUrlGithub),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDbtCloudProjectRepositoryExists("dbt_cloud_project_repository.test_project_repository"),
				),
			},
			// IMPORT
			{
				ResourceName:            "dbt_cloud_project_repository.test_project_repository",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
			// EMPTY
			{
				Config: testAccDbtCloudProjectRepositoryResourceEmptyConfig(projectName, repoUrlGithub),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDbtCloudProjectRepositoryEmpty("dbt_cloud_project.test_project"),
				),
			},
		},
	})
}

func testAccDbtCloudProjectRepositoryResourceBasicConfig(projectName, repoUrlGithub string) string {
	return fmt.Sprintf(`
resource "dbt_cloud_project" "test_project" {
  name        = "%s"
}

resource "dbt_cloud_repository" "test_repository" {
  remote_url = "%s"
  project_id = dbt_cloud_project.test_project.id
  depends_on = [dbt_cloud_project.test_project]
}

resource "dbt_cloud_project_repository" "test_project_repository" {
  project_id = dbt_cloud_project.test_project.id
  repository_id = dbt_cloud_repository.test_repository.repository_id
}
`, projectName, repoUrlGithub)
}

func testAccDbtCloudProjectRepositoryResourceEmptyConfig(projectName, repoUrlGithub string) string {
	return fmt.Sprintf(`
resource "dbt_cloud_project" "test_project" {
  name        = "%s"
}

resource "dbt_cloud_repository" "test_repository" {
  remote_url = "%s"
  project_id = dbt_cloud_project.test_project.id
  depends_on = [dbt_cloud_project.test_project]
}
`, projectName, repoUrlGithub)
}

func testAccCheckDbtCloudProjectRepositoryExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}
		apiClient := testAccProvider.Meta().(*dbt_cloud.Client)
		projectId := strings.Split(rs.Primary.ID, dbt_cloud.ID_DELIMITER)[0]
		project, err := apiClient.GetProject(projectId)
		if err != nil {
			return fmt.Errorf("Can't get project")
		}
		if project.RepositoryID == nil {
			return fmt.Errorf("error fetching item with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccCheckDbtCloudProjectRepositoryEmpty(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}
		apiClient := testAccProvider.Meta().(*dbt_cloud.Client)
		project, err := apiClient.GetProject(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Can't get project")
		}
		if project.RepositoryID != nil {
			return fmt.Errorf("error fetching item with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccCheckDbtCloudProjectRepositoryDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*dbt_cloud.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "dbt_cloud_project_repository" {
			continue
		}
		projectId := strings.Split(rs.Primary.ID, dbt_cloud.ID_DELIMITER)[0]
		project, err := apiClient.GetProject(projectId)
		if project != nil {
			return fmt.Errorf("Project still exists")
		}
		notFoundErr := "not found"
		expectedErr := regexp.MustCompile(notFoundErr)
		if !expectedErr.Match([]byte(err.Error())) {
			return fmt.Errorf("expected %s, got %s", notFoundErr, err)
		}
	}

	return nil
}
