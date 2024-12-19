package job_test

import (
	"fmt"
	"github.com/dbt-labs/terraform-provider-dbtcloud/pkg/framework/acctest_config"
	"testing"

	"github.com/dbt-labs/terraform-provider-dbtcloud/pkg/framework/acctest_helper"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDbtCloudJobsDataSource(t *testing.T) {

	randomJobName := acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)
	randomJobName2 := acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)

	config := jobs(randomJobName, randomJobName2)

	check := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrSet("data.dbtcloud_jobs.test", "project_id"),
		resource.TestCheckResourceAttr("data.dbtcloud_jobs.test", "jobs.#", "2"),
		resource.TestCheckResourceAttrSet("data.dbtcloud_jobs.test", "jobs.0.id"),
		resource.TestCheckResourceAttrSet("data.dbtcloud_jobs.test", "jobs.0.name"),
		resource.TestCheckResourceAttrSet("data.dbtcloud_jobs.test", "jobs.1.id"),
		resource.TestCheckResourceAttrSet("data.dbtcloud_jobs.test", "jobs.1.name"),

		resource.TestCheckResourceAttrSet("data.dbtcloud_jobs.test_env", "environment_id"),
		resource.TestCheckResourceAttr("data.dbtcloud_jobs.test_env", "jobs.#", "1"),
		resource.TestCheckResourceAttrSet("data.dbtcloud_jobs.test_env", "jobs.0.id"),
		resource.TestCheckResourceAttrSet(
			"data.dbtcloud_jobs.test_env",
			"jobs.0.project_id",
		),
		resource.TestCheckResourceAttrSet("data.dbtcloud_jobs.test_env", "environment_id"),
		resource.TestCheckResourceAttr("data.dbtcloud_jobs.test_env", "jobs.0.name", randomJobName),
		resource.TestCheckResourceAttr(
			"data.dbtcloud_jobs.test_env",
			"jobs.0.environment.deployment_type",
			"production",
		),
		resource.TestCheckResourceAttr(
			"data.dbtcloud_jobs.test_env",
			"jobs.0.execution.timeout_seconds",
			"180",
		),
		resource.TestCheckResourceAttr(
			"data.dbtcloud_jobs.test_env",
			"jobs.0.triggers_on_draft_pr",
			"false",
		),
		resource.TestCheckResourceAttr(
			"data.dbtcloud_jobs.test_env",
			"jobs.0.job_completion_trigger_condition.condition.statuses.0",
			"success",
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

func jobs(jobName string, jobName2 string) string {
	return fmt.Sprintf(`
    resource "dbtcloud_project" "test_project" {
        name = "jobs_test_project"
    }

    resource "dbtcloud_environment" "test_environment" {
        project_id = dbtcloud_project.test_project.id
        name = "job_test_env"
        dbt_version = "%s"
        type = "deployment"
		deployment_type = "production"
    }

    resource "dbtcloud_environment" "test_environment2" {
        project_id = dbtcloud_project.test_project.id
        name = "job_test_env2"
        dbt_version = "%s"
        type = "deployment"
    }

    resource "dbtcloud_job" "test_job" {
        name = "%s"
        project_id = dbtcloud_project.test_project.id
        environment_id = dbtcloud_environment.test_environment.environment_id
        execute_steps = [
            "dbt run"
        ]
        triggers = {
          "github_webhook" : false,
          "schedule" : false,
          "git_provider_webhook": false
        }
        timeout_seconds = 180
		job_completion_trigger_condition {
			job_id = dbtcloud_job.test_job2.id
			project_id = dbtcloud_project.test_project.id
			statuses = ["success"]
		}
    }

    resource "dbtcloud_job" "test_job2" {
        name = "%s"
        project_id = dbtcloud_project.test_project.id
        environment_id = dbtcloud_environment.test_environment2.environment_id
        execute_steps = [
            "dbt run"
        ]
        triggers = {
          "github_webhook" : false,
          "schedule" : false,
          "git_provider_webhook": false
        }
        timeout_seconds = 1800
    }	

    data "dbtcloud_jobs" "test" {
        project_id = dbtcloud_project.test_project.id
		depends_on = [
			dbtcloud_job.test_job,
			dbtcloud_job.test_job2,
			]
	}
			
	data "dbtcloud_jobs" "test_env" {
        environment_id = dbtcloud_environment.test_environment.environment_id
		  depends_on = [
			dbtcloud_job.test_job,
			dbtcloud_job.test_job2,
		]
    }
    `, acctest_config.AcceptanceTestConfig.DbtCloudVersion, acctest_config.AcceptanceTestConfig.DbtCloudVersion, jobName, jobName2)
}
