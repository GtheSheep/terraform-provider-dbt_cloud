package environment_test

import (
	"fmt"
	"github.com/dbt-labs/terraform-provider-dbtcloud/pkg/framework/acctest_config"
	"testing"

	"github.com/dbt-labs/terraform-provider-dbtcloud/pkg/framework/acctest_helper"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDbtCloudEnvironmentsDataSource(t *testing.T) {

	randomProjectName1 := acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)
	randomProjectName2 := acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)
	randomEnvironmentName1 := acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)
	randomEnvironmentName2 := acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)

	config := environments(
		randomProjectName1,
		randomProjectName2,
		randomEnvironmentName1,
		randomEnvironmentName2,
	)

	check := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrSet(
			"data.dbtcloud_environments.test_all",
			"environments.1.%",
		),
		resource.TestCheckResourceAttr(
			"data.dbtcloud_environments.test_one_project",
			"environments.#",
			"1",
		),
		resource.TestCheckResourceAttr(
			"data.dbtcloud_environments.test_one_project",
			"environments.0.name",
			randomEnvironmentName2,
		),
		resource.TestCheckResourceAttrSet(
			"data.dbtcloud_environments.test_one_project",
			"environments.0.environment_id",
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

func environments(
	randomProjectName1,
	randomProjectName2,
	randomEnvironmentName1,
	randomEnvironmentName2 string) string {
	return fmt.Sprintf(`
    resource "dbtcloud_project" "test_project1" {
        name = "%s"
    }

	resource "dbtcloud_project" "test_project2" {
        name = "%s"
    }

    resource "dbtcloud_environment" "test_environment1" {
        project_id = dbtcloud_project.test_project1.id
        name = "%s"
        dbt_version = "%s"
        type = "deployment"
        use_custom_branch = true
        custom_branch = "customBranchName"
		enable_model_query_history = true
    }

	resource "dbtcloud_environment" "test_environment2" {
        project_id = dbtcloud_project.test_project2.id
        name = "%s"
        dbt_version = "%s"
        type = "deployment"
        use_custom_branch = false
    }

    data "dbtcloud_environments" "test_all" {
		depends_on = [dbtcloud_environment.test_environment1, dbtcloud_environment.test_environment2]
    }

	data "dbtcloud_environments" "test_one_project" {
        project_id = dbtcloud_project.test_project2.id
		depends_on = [dbtcloud_environment.test_environment1, dbtcloud_environment.test_environment2]
    }
    `,
		randomProjectName1,
		randomProjectName2,
		randomEnvironmentName1,
		acctest_config.AcceptanceTestConfig.DbtCloudVersion,
		randomEnvironmentName2,
		acctest_config.AcceptanceTestConfig.DbtCloudVersion,
	)
}
