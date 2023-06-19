---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "dbtcloud_privatelink_endpoint Data Source - dbtcloud"
subcategory: ""
description: |-
  
---

# dbtcloud_privatelink_endpoint (Data Source)



## Example Usage

```terraform
// use dbt_cloud_privatelink_endpoint instead of dbtcloud_privatelink_endpoint for the legacy resource names
// legacy names will be removed from 0.3 onwards

data "dbtcloud_privatelink_endpoint" "test_with_name" {
  name = "My Endpoint Name"
}

data "dbtcloud_privatelink_endpoint" "test_with_url" {
  private_link_endpoint_url = "abc.privatelink.def.com"

}
// in case multiple endpoints have the same name or URL
data "dbtcloud_privatelink_endpoint" "test_with_name_and_url" {
  name = "My Endpoint Name"
  private_link_endpoint_url = "abc.privatelink.def.com"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `name` (String) Given descriptive name for the PrivateLink Endpoint (name and/or private_link_endpoint_url need to be provided to return data for the datasource)
- `private_link_endpoint_url` (String) The URL of the PrivateLink Endpoint (private_link_endpoint_url and/or name need to be provided to return data for the datasource)

### Read-Only

- `cidr_range` (String) The CIDR range of the PrivateLink Endpoint
- `id` (String) The internal ID of the PrivateLink Endpoint
- `state` (Number) PrivatelinkEndpoint state should be 1 = active, as 2 = deleted
- `type` (String) Type of the PrivateLink Endpoint

