package user_test

import (
	"github.com/dbt-labs/terraform-provider-dbtcloud/pkg/framework/acctest_config"
	"testing"

	"github.com/dbt-labs/terraform-provider-dbtcloud/pkg/framework/acctest_helper"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDbtCloudUsersDataSource(t *testing.T) {

	userEmail := acctest_config.AcceptanceTestConfig.DbtCloudUserEmail

	_ = userEmail
	config := users()

	check := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrSet("data.dbtcloud_users.all", "users.0.email"),
		resource.TestCheckResourceAttrSet("data.dbtcloud_users.all", "users.0.id"),
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

func users() string {
	return `
data "dbtcloud_users" "all" {
}
`
}
