---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "dbt-cloud_environment Resource - terraform-provider-dbt-cloud"
subcategory: ""
description: |-
  
---

# dbt-cloud_environment (Resource)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `dbt_version` (String) Version number of dbt to use in this environment
- `name` (String) Environment name
- `project_id` (Number) Project ID to create the environment in
- `type` (String) The type of environment (must be either development or deployment)

### Optional

- `credential_id` (Number) Credential ID to create the environment with
- `custom_branch` (String) Which custom branch to use in this environment
- `is_active` (Boolean) Whether the environment is active
- `use_custom_branch` (Boolean) Whether to use a custom git branch in this environment

### Read-Only

- `environment_id` (Number) Environment ID within the project
- `id` (String) The ID of this resource.

