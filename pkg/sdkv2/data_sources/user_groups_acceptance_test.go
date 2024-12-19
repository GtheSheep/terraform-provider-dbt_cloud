package data_sources_test

import (
	"fmt"
	"github.com/dbt-labs/terraform-provider-dbtcloud/pkg/framework/acctest_config"
	"testing"

	"github.com/dbt-labs/terraform-provider-dbtcloud/pkg/framework/acctest_helper"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDbtCloudUserGroupsDataSource(t *testing.T) {

	userID := acctest_config.AcceptanceTestConfig.DbtCloudUserId
	config := userGroups(userID)

	check := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrSet("data.dbtcloud_user_groups.test", "user_id"),
		resource.TestCheckResourceAttrSet("data.dbtcloud_user_groups.test", "group_ids.0"),
		resource.TestCheckResourceAttrSet("data.dbtcloud_user_groups.test", "group_ids.1"),
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

func userGroups(userID int) string {
	return fmt.Sprintf(`
    data "dbtcloud_user_groups" "test" {
        user_id = %d
    }
    `, userID)
}
