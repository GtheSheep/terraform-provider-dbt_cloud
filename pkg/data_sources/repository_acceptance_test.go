package data_sources_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDbtCloudRepositoryDataSource(t *testing.T) {

	randomProjectName := acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)
	repoUrl := "git@github.com:GtheSheep/terraform-provider-dbt_cloud.git"

	config := repository(randomProjectName, repoUrl)

	check := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr("data.dbt_cloud_repository.test", "remote_url", repoUrl),
		resource.TestCheckResourceAttrSet("data.dbt_cloud_repository.test", "repository_id"),
		resource.TestCheckResourceAttrSet("data.dbt_cloud_repository.test", "project_id"),
		resource.TestCheckResourceAttrSet("data.dbt_cloud_repository.test", "is_active"),
	)

	resource.ParallelTest(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check:  check,
			},
		},
	})
}

func repository(projectName, repositoryUrl string) string {
	return fmt.Sprintf(`
    resource "dbt_cloud_project" "test_project" {
        name = "%s"
    }

    resource "dbt_cloud_repository" "test_repository" {
        project_id = dbt_cloud_project.test_project.id
        remote_url = "%s"
        is_active = true
        depends_on = [
            dbt_cloud_project.test_project
        ]
    }

    data "dbt_cloud_repository" "test" {
        project_id = dbt_cloud_project.test_project.id
        repository_id = dbt_cloud_repository.test_repository.repository_id
    }
    `, projectName, repositoryUrl)
}
