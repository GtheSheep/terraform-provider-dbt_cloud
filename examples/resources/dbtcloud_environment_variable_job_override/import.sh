# Import using a project ID, a job ID and the environment variable override ID
terraform import dbtcloud_environment_variable_job_override.test_environment_variable_job_override "project_id:job_id:environment_variable_override_id"
terraform import dbtcloud_environment_variable_job_override.test_environment_variable_job_override 12345:678:123456