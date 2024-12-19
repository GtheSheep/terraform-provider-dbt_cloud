package user_test

import (
	"fmt"
	"github.com/dbt-labs/terraform-provider-dbtcloud/pkg/framework/acctest_config"
	"testing"

	"github.com/dbt-labs/terraform-provider-dbtcloud/pkg/framework/acctest_helper"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDbtCloudUserDataSource(t *testing.T) {

	userEmail := acctest_config.AcceptanceTestConfig.DbtCloudUserEmail

	config := user(userEmail)

	check := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr("data.dbtcloud_user.test_user_read", "email", userEmail),
		resource.TestCheckResourceAttrSet("data.dbtcloud_user.test_user_read", "id"),
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

func user(userEmail string) string {
	return fmt.Sprintf(`
data "dbtcloud_user" "test_user_read" {
    email = "%s"
}
`, userEmail)
}
