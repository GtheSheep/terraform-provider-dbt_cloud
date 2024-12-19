package acctest_config

import (
	"log"
	"os"
	"strconv"
)

var AcceptanceTestConfig = buildAcctestConfig()

func buildAcctestConfig() AcctestConfig {
	return AcctestConfig{
		DbtCloudAccountId:           determineIntValue("DBT_CLOUD_ACCOUNT_ID", 1, 1),
		DbtCloudServiceToken:        os.Getenv("DBT_CLOUD_TOKEN"),
		DbtCloudPersonalAccessToken: os.Getenv("DBT_CLOUD_PERSONAL_ACCESS_TOKEN"),
		DbtCloudHostUrl:             determineStringValue("DBT_CLOUD_HOST_URL", "", ""),
		DbtCloudVersion:             "latest",

		DbtCloudUserId: determineIntValue(
			"ACC_TEST_DBT_CLOUD_USER_ID",
			1,
			54461,
		),
		DbtCloudUserEmail: determineStringValue(
			"ACC_TEST_DBT_CLOUD_USER_EMAIL",
			"d"+"ev@"+"db"+"tla"+"bs.c"+"om",
			"beno"+"it"+".per"+"igaud"+"@"+"fisht"+"ownanalytics"+"."+"com",
		),
		DbtCloudGroupIds: determineStringValue(
			"ACC_TEST_DBT_CLOUD_GROUP_IDS",
			"1,2,3",
			"531585,531584,531583",
		),
		AzureDevOpsProjectName: determineStringValue(
			"ACC_TEST_AZURE_DEVOPS_PROJECT_NAME",
			"dbt-cloud-ado-project",
			"dbt-cloud-ado-project",
		),
		GitHubRepoUrl: determineStringValue(
			"ACC_TEST_GITHUB_REPO_URL",
			"git://github.com/dbt-labs/jaffle_shop.git",
			"git://github.com/dbt-labs/jaffle_shop.git",
		),
		GitHubAppInstallationId: determineIntValue(
			"ACC_TEST_GITHUB_APP_INSTALLATION_ID",
			28374841,
			28374841,
		),
	}
}

type AcctestConfig struct {
	DbtCloudAccountId           int
	DbtCloudServiceToken        string
	DbtCloudPersonalAccessToken string
	DbtCloudHostUrl             string
	DbtCloudVersion             string

	DbtCloudUserId          int
	DbtCloudUserEmail       string
	DbtCloudGroupIds        string
	AzureDevOpsProjectName  string
	GitHubRepoUrl           string
	GitHubAppInstallationId int
}

func IsDbtCloudPR() bool {
	return os.Getenv("DBT_CLOUD_ACCOUNT_ID") == "1"
}

func IsCI() bool {
	return os.Getenv("CI") != ""
}

func determineStringValue(envVarKey string, dbtCloudPRValue string, ciValue string) string {
	val := os.Getenv(envVarKey)
	if val != "" {
		return val
	} else if IsDbtCloudPR() {
		return dbtCloudPRValue
	} else if IsCI() {
		return ciValue
	} else {
		log.Printf("Unable to determine %s value, tests may fail", envVarKey)
		return ""
	}
}

func determineIntValue(envVarKey string, dbtCloudPRValue int, ciValue int) int {
	val := os.Getenv(envVarKey)
	if val != "" {
		intVal, err := strconv.Atoi(val)
		if err != nil {
			log.Fatalf("Unable to determine %s value for test: %v", envVarKey, err)
		}
		return intVal
	} else if IsDbtCloudPR() {
		return dbtCloudPRValue
	} else if IsCI() {
		return ciValue
	} else {
		log.Printf("Unable to determine %s value, tests may fail", envVarKey)
		return -1
	}
}
