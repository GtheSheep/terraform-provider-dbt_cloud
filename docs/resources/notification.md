---
page_title: "dbtcloud_notification Resource - dbtcloud"
subcategory: ""
description: |-
  
---

# dbtcloud_notification (Resource)




## Example Usage

```terraform
// dbt Cloud allows us to create internal and external notifications

// an internal notification will send emails to the user mentioned in `user_id`
resource "dbtcloud_notification" "my_internal_notification" {
	// user_id is the internal ID of a given user in dbt Cloud
	user_id           = 100
	on_success        = [dbtcloud_job.my_job.id]
	on_failure        = [12345]
	// the Type 1 is used for internal notifications
	notification_type = 1
}

// we can also send "external" email notifications to emails to related to dbt Cloud users
resource "dbtcloud_notification" "my_external_notification" {
	// we still need the ID of a user in dbt Cloud even though it is not used for sending notifications
	user_id           = 100
	on_failure        = [23456, 56788]
	on_cancel         = [dbtcloud_job.my_other_job.id]
	// the Type 4 is used for external notifications
	notification_type = 4
	// the external_email is the email address that will receive the notification
	external_email    = "my_email@mail.com"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `user_id` (Number) Internal dbt Cloud User ID. Must be the user_id for an existing user even if the notification is an external one

### Optional

- `external_email` (String) The external email to receive the notification
- `notification_type` (Number) Type of notification (1 = dbt Cloud user email (default): does not require an external_email ; 4 = external email: requires setting an external_email)
- `on_cancel` (Set of Number) List of job IDs to trigger the webhook on cancel
- `on_failure` (Set of Number) List of job IDs to trigger the webhook on failure
- `on_success` (Set of Number) List of job IDs to trigger the webhook on success
- `state` (Number) State of the notification (1 = active (default), 2 = inactive)

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
# Import using a notification ID
terraform import dbtcloud_notification.my_notification "notification_id"
terraform import dbtcloud_notification.my_notification 12345
```