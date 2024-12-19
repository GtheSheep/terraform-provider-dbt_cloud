package partial_notification_test

import (
	"fmt"
	"github.com/dbt-labs/terraform-provider-dbtcloud/pkg/framework/acctest_config"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/dbt-labs/terraform-provider-dbtcloud/pkg/framework/acctest_helper"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccDbtCloudPartialNotificationResource(t *testing.T) {

	if acctest_config.IsDbtCloudPR() {
		t.Skip("Skipping notifications in dbt Cloud CI for now")
	}

	userID := acctest_config.AcceptanceTestConfig.DbtCloudUserId

	currentTime := time.Now().Unix()
	notificationEmail := fmt.Sprintf("%d-partial-resource@nomail.com", currentTime)

	projectName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest_helper.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest_helper.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDbtCloudPartialNotificationResourceCreatePartialNotifications(
					projectName,
					userID,
					notificationEmail,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDbtCloudPartialNotificationExists(
						"dbtcloud_partial_notification.test_notification_internal",
					),
					resource.TestCheckResourceAttr(
						"dbtcloud_partial_notification.test_notification_internal",
						"notification_type",
						"1",
					),
					resource.TestCheckResourceAttrSet(
						"dbtcloud_partial_notification.test_notification_internal",
						"on_success.0",
					),
					resource.TestCheckResourceAttrSet(
						"dbtcloud_partial_notification.test_notification_internal",
						"on_cancel.0",
					),
					resource.TestCheckResourceAttrSet(
						"dbtcloud_partial_notification.test_notification_internal",
						"on_cancel.1",
					),
					resource.TestCheckResourceAttrSet(
						"dbtcloud_partial_notification.test_notification_internal",
						"on_failure.0",
					),

					testAccCheckDbtCloudPartialNotificationExists(
						"dbtcloud_partial_notification.test_notification_external",
					),
					resource.TestCheckResourceAttr(
						"dbtcloud_partial_notification.test_notification_external",
						"notification_type",
						"4",
					),
					resource.TestCheckResourceAttrSet(
						"dbtcloud_partial_notification.test_notification_external",
						"on_warning.0",
					),
					resource.TestCheckResourceAttrSet(
						"dbtcloud_partial_notification.test_notification_external",
						"on_failure.0",
					),
					resource.TestCheckResourceAttr(
						"dbtcloud_partial_notification.test_notification_external",
						"external_email",
						notificationEmail,
					),
				),
			},
			// MODIFY
			{
				Config: testAccDbtCloudPartialNotificationResourceModifyPartialNotifications(
					projectName,
					userID,
					notificationEmail,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDbtCloudPartialNotificationExists(
						"dbtcloud_partial_notification.test_notification_internal",
					),
					resource.TestCheckNoResourceAttr(
						"dbtcloud_partial_notification.test_notification_internal",
						"on_cancel.0",
					),
					resource.TestCheckResourceAttrSet(
						"dbtcloud_partial_notification.test_notification_internal",
						"on_warning.0",
					),
					resource.TestCheckResourceAttrSet(
						"dbtcloud_partial_notification.test_notification_internal",
						"on_warning.1",
					),

					testAccCheckDbtCloudPartialNotificationExists(
						"dbtcloud_partial_notification.test_notification_external",
					),
					resource.TestCheckResourceAttrSet(
						"dbtcloud_partial_notification.test_notification_external",
						"on_cancel.0",
					),
					resource.TestCheckResourceAttrSet(
						"dbtcloud_partial_notification.test_notification_external2",
						"on_success.0",
					),
					resource.TestCheckResourceAttr(
						"dbtcloud_partial_notification.test_notification_external2",
						"on_cancel.#",
						"0",
					),
				),
			},
		},
	})
}

func testAccDbtCloudPartialNotificationResourceBasicConfig(projectName string) string {
	return fmt.Sprintf(`
resource "dbtcloud_project" "test_notification_project" {
	name = "%s"
}
	
resource "dbtcloud_environment" "test_notification_environment" {
	project_id  = dbtcloud_project.test_notification_project.id
	name        = "Test Env PartialNotification"
	dbt_version = "%s"
	type        = "development"
}
	
resource "dbtcloud_job" "test_notification_job_1" {
	name           = "Job 1 TF"
	project_id     = dbtcloud_project.test_notification_project.id
	environment_id = dbtcloud_environment.test_notification_environment.environment_id
	execute_steps = [
		"dbt compile"
	]
	triggers = {
		"github_webhook" : false,
		"git_provider_webhook" : false,
		"schedule" : false,
	}
}
	
resource "dbtcloud_job" "test_notification_job_2" {
	name           = "Job 2 TF"
	project_id     = dbtcloud_project.test_notification_project.id
	environment_id = dbtcloud_environment.test_notification_environment.environment_id
	execute_steps = [
		"dbt test"
	]
	triggers = {
		"github_webhook" : false,
		"git_provider_webhook" : false,
		"schedule" : false,
	}
}
`, projectName, acctest_config.AcceptanceTestConfig.DbtCloudVersion)
}

func testAccDbtCloudPartialNotificationResourceCreatePartialNotifications(
	projectName string,
	userID int,
	notificationEmail string,
) string {

	notificationsConfig := fmt.Sprintf(`
resource "dbtcloud_partial_notification" "test_notification_internal" {
	user_id           = %d
	on_success        = [dbtcloud_job.test_notification_job_1.id]
	on_failure        = [dbtcloud_job.test_notification_job_2.id]
	on_cancel         = [dbtcloud_job.test_notification_job_1.id, dbtcloud_job.test_notification_job_2.id]
	notification_type = 1
}
	
resource "dbtcloud_partial_notification" "test_notification_external" {
	user_id           = %d
	on_warning        = [dbtcloud_job.test_notification_job_1.id]
	on_failure        = [dbtcloud_job.test_notification_job_2.id]
	notification_type = 4
	external_email    = "%s"
}
`, userID, userID, notificationEmail)
	return testAccDbtCloudPartialNotificationResourceBasicConfig(
		projectName,
	) + "\n" + notificationsConfig
}

func testAccDbtCloudPartialNotificationResourceModifyPartialNotifications(
	projectName string,
	userID int,
	notificationEmail string,
) string {

	notificationsConfig := fmt.Sprintf(`
resource "dbtcloud_partial_notification" "test_notification_internal" {
	user_id           = %d
	on_success        = [dbtcloud_job.test_notification_job_1.id]
	on_failure        = [dbtcloud_job.test_notification_job_2.id]
	on_cancel         = []
	on_warning        = [dbtcloud_job.test_notification_job_1.id, dbtcloud_job.test_notification_job_2.id]
	notification_type = 1
}
	
resource "dbtcloud_partial_notification" "test_notification_external" {
	user_id           = %d
	on_failure        = [dbtcloud_job.test_notification_job_2.id]
	on_cancel         = [dbtcloud_job.test_notification_job_1.id]
	notification_type = 4
	external_email    = "%s"
}

resource "dbtcloud_partial_notification" "test_notification_external2" {
	user_id           = %d
	on_success        = [dbtcloud_job.test_notification_job_1.id, dbtcloud_job.test_notification_job_2.id]
	notification_type = 4
	external_email    = "%s"
	depends_on        = [dbtcloud_partial_notification.test_notification_external]
}
`, userID, userID, notificationEmail, userID, notificationEmail)
	return testAccDbtCloudPartialNotificationResourceBasicConfig(
		projectName,
	) + "\n" + notificationsConfig
}

func testAccCheckDbtCloudPartialNotificationExists(resource string) resource.TestCheckFunc {
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

		_, err = apiClient.GetNotification(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error fetching item with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccCheckDbtCloudPartialNotificationDestroy(s *terraform.State) error {

	apiClient, err := acctest_helper.SharedClient()
	if err != nil {
		return fmt.Errorf("Issue getting the client")
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "dbtcloud_partial_notification" {
			continue
		}
		_, err := apiClient.GetNotification(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("PartialNotification still exists")
		}
		notFoundErr := "resource-not-found"
		expectedErr := regexp.MustCompile(notFoundErr)
		if !expectedErr.Match([]byte(err.Error())) {
			return fmt.Errorf("expected %s, got %s", notFoundErr, err)
		}
	}

	return nil
}
