package resources_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/gthesheep/terraform-provider-dbt_cloud/pkg/dbt_cloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDbtCloudProjectResource(t *testing.T) {

	projectName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	projectName2 := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDbtCloudProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDbtCloudProjectResourceBasicConfig(projectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDbtCloudProjectExists("dbt_cloud_project.test_project"),
					resource.TestCheckResourceAttr("dbt_cloud_project.test_project", "name", projectName),
				),
			},
			// RENAME
			{
				Config: testAccDbtCloudProjectResourceBasicConfig(projectName2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDbtCloudProjectExists("dbt_cloud_project.test_project"),
					resource.TestCheckResourceAttr("dbt_cloud_project.test_project", "name", projectName2),
				),
			},
			// MODIFY
			{
				Config: testAccDbtCloudProjectResourceFullConfig(projectName2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDbtCloudProjectExists("dbt_cloud_project.test_project"),
					resource.TestCheckResourceAttr("dbt_cloud_project.test_project", "name", projectName2),
					resource.TestCheckResourceAttr("dbt_cloud_project.test_project", "dbt_project_subdirectory", "/project/subdirectory_where/dbt-is"),
				),
			},
			// IMPORT
			{
				ResourceName:            "dbt_cloud_project.test_project",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func testAccDbtCloudProjectResourceBasicConfig(projectName string) string {
	return fmt.Sprintf(`
resource "dbt_cloud_project" "test_project" {
  name        = "%s"
}
`, projectName)
}

func testAccDbtCloudProjectResourceFullConfig(projectName string) string {
	return fmt.Sprintf(`
resource "dbt_cloud_project" "test_project" {
  name        = "%s"
  dbt_project_subdirectory = "/project/subdirectory_where/dbt-is"
}
`, projectName)
}

func testAccCheckDbtCloudProjectExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}
		apiClient := testAccProvider.Meta().(*dbt_cloud.Client)
		_, err := apiClient.GetProject(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error fetching item with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccCheckDbtCloudProjectDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*dbt_cloud.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "dbt_cloud_project" {
			continue
		}
		_, err := apiClient.GetProject(rs.Primary.ID)
		if err == nil {
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
