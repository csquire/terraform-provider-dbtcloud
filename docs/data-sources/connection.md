---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "dbtcloud_connection Data Source - dbtcloud"
subcategory: ""
description: |-
  
---

# dbtcloud_connection (Data Source)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `connection_id` (Number) ID for the connection
- `project_id` (Number) Project ID to create the connection in

### Read-Only

- `account` (String) Account for the connection
- `allow_keep_alive` (Boolean) Flag for whether or not to use the keep session alive parameter in the connection
- `allow_sso` (Boolean) Flag for whether or not to use SSO for the connection
- `database` (String) Database name for the connection
- `id` (String) The ID of this resource.
- `is_active` (Boolean) Whether the connection is active
- `name` (String) Connection name
- `private_link_endpoint_id` (String) The ID of the PrivateLink connection
- `role` (String) Role name for the connection
- `type` (String) Connection type
- `warehouse` (String) Warehouse name for the connection

