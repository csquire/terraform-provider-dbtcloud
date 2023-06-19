// use dbt_cloud_environment instead of dbtcloud_environment for the legacy resource names
// legacy names will be removed from 0.3 onwards

resource "dbtcloud_environment" "test_environment" {
  dbt_version   = "1.0.1"
  name          = "test"
  project_id    = data.dbt_cloud_project.test_project.project_id
  type          = "deployment"
  credential_id = dbt_cloud_snowflake_credential.new_credential.credential_id
}