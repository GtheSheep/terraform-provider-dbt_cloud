# terraform-provider-dbtcloud

## Terraform Provider for dbt Cloud

This repo was originally created by a dbt community member, Gary James [[GtheSheep](https://github.com/GtheSheep)]

> If you are still using the GtheSheep/dbt-cloud source, see [Upgrading from the community Provider](UPGRADING_PROVIDER.md) to upgrade to the latest version.

## Scope

Provide the ability to manage dbt Cloud projects and account settings via Terraform resources.
Data sources are also available for most resources.

In order to use this provider, add the following to your Terraform providers
setup, with the latest version number.

```terraform
terraform {
  required_providers {
    dbtcloud = {
      source  = "dbt-labs/dbtcloud"
      version = "<version>"
    }
  }
}
```

## Authentication

If you want to explicitly set the authentication variables on the provider, you
can do so as below, though likely via a `variables.tf` file or config in your
CI-CD pipeline to keep these credentials safe.

```terraform
provider "dbtcloud" {
  // required
  account_id = ...
  token      = "..."
  // optional - defaults to the US Multi Tenant URL "https://cloud.getdbt.com/api"
  host_url   = "..."
}
```

You can also set them via environment variables:  
`DBT_CLOUD_ACCOUNT_ID` for the `account_id`.  
`DBT_CLOUD_TOKEN` for the `token`.  
`DBT_CLOUD_HOST_URL` (Optional) for the `host_url`.

## Getting started and Examples

The provider documentation is directly available [on the Terraform Registry](https://registry.terraform.io/providers/dbt-labs/dbtcloud/latest/docs).

- Under [Guides](https://registry.terraform.io/providers/dbt-labs/dbtcloud/latest/docs/guides/1_getting_started), you will find a simple example of how to use the provider
- Each resource ([example for jobs](https://registry.terraform.io/providers/dbt-labs/dbtcloud/latest/docs/resources/job)) has some usage examples and contains the list of parameters available

## Importing existing dbt Cloud configuration

The CLI [dbtcloud-terraforming](https://github.com/dbt-labs/dbtcloud-terraforming) can be used to generate the Terraform configuration and import statements based on your existing dbt Cloud configuration.

## Running Acceptance Tests

Acceptance tests are executed by running the `make test-acceptance` command.

For the acceptance tests to work locally, the following environment variables must be set to appropriate values
for a dbt Cloud account the tests can interact with. All dbt Cloud resources referenced by the environment variables
(e.g. user id, email address, and group ids) must exist in the dbt Cloud account.
```
DBT_CLOUD_ACCOUNT_ID=1234
DBT_CLOUD_HOST_URL=https://<host>/api
DBT_CLOUD_TOKEN=<api_token>
DBT_CLOUD_PERSONAL_ACCESS_TOKEN=<api_token>
ACC_TEST_DBT_CLOUD_USER_EMAIL=<email_address>
ACC_TEST_DBT_CLOUD_USER_ID=4321
ACC_TEST_DBT_CLOUD_GROUP_IDS=1,2,3
ACC_TEST_AZURE_DEVOPS_PROJECT_NAME=test-project
ACC_TEST_GITHUB_REPO_URL=git@github.com:<github-org>/<dbt-project>.git
ACC_TEST_GITHUB_APP_INSTALLATION_ID=1234
```

To assist with setting the environment variables, the `.env.example` file can be copied to `.env` and the values updated.  
The variables can then be loaded into the environment by running `source .env` on Mac or Linux.

**A note on the Repository Acceptance Tests**  
The Repository Acceptance Tests require a GitHub repository to be set up and the dbt Cloud GitHub App installed.

`ACC_TEST_GITHUB_REPO_URL` must be set to the SSH URL of a repository

`ACC_TEST_GITHUB_APP_INSTALLATION_ID` must be set to the installation ID of the GitHub App.  
The installation ID can be found by navigating to `Settings` -> `Applications`, 
and clicking `Configure` on the dbt Cloud GitHub App. The installation ID can be found in the url, for example,
`https://github.com/settings/installations/<installation_id>`

## Acknowledgement

Thanks to Gary James [[GtheSheep](https://github.com/GtheSheep)], for all the effort put in creating this provider originally
and for being a great dbt community member!
