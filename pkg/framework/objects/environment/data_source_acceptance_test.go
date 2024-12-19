package environment_test

import (
	"fmt"
	"github.com/dbt-labs/terraform-provider-dbtcloud/pkg/framework/acctest_config"
	"testing"

	"github.com/dbt-labs/terraform-provider-dbtcloud/pkg/framework/acctest_helper"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDbtCloudEnvironmentDataSource(t *testing.T) {

	randomProjectName := acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)
	randomEnvironmentName := acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)

	config := environment(randomProjectName, randomEnvironmentName)

	check := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(
			"data.dbtcloud_environment.test",
			"name",
			randomEnvironmentName,
		),
		resource.TestCheckResourceAttrSet("data.dbtcloud_environment.test", "environment_id"),
		resource.TestCheckResourceAttrSet("data.dbtcloud_environment.test", "project_id"),
		resource.TestCheckResourceAttrSet("data.dbtcloud_environment.test", "dbt_version"),
		resource.TestCheckResourceAttrSet("data.dbtcloud_environment.test", "type"),
		resource.TestCheckResourceAttrSet("data.dbtcloud_environment.test", "use_custom_branch"),
		resource.TestCheckResourceAttr(
			"data.dbtcloud_environment.test",
			"custom_branch",
			"customBranchName",
		),
		resource.TestCheckResourceAttr(
			"data.dbtcloud_environment.test",
			"enable_model_query_history",
			"true",
		),
	)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest_helper.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check:  check,
			},
		},
	})
}

func environment(projectName, environmentName string) string {
	return fmt.Sprintf(`
    resource "dbtcloud_project" "test_project" {
        name = "%s"
    }

    resource "dbtcloud_environment" "test_environment" {
        project_id = dbtcloud_project.test_project.id
        name = "%s"
        dbt_version = "%s"
        type = "development"
        use_custom_branch = true
        custom_branch = "customBranchName"
		enable_model_query_history = true
    }

    data "dbtcloud_environment" "test" {
        project_id = dbtcloud_project.test_project.id
        environment_id = dbtcloud_environment.test_environment.environment_id
    }
    `, projectName, environmentName, acctest_config.AcceptanceTestConfig.DbtCloudVersion)
}
