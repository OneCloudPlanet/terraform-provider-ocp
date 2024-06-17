---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "ocp_cluster_version Data Source - terraform-provider-ocp"
subcategory: ""
description: |-
  List available cluster addons for Kubernetes Cluster.
---

# ocp_cluster_version

List available cluster addons for Kubernetes Cluster.

## Example Usage

```hcl
data "ocp_cluster_addons" "list_addons" {
  filter {
    name = "ingress-nginx"
  }
}
```

## Argument Reference

- `filter` - (Optional) Values to filter available addons:
    + `name` - (String) Filter by addons name.
    + `version` - (String) Filter by addons version

## Attributes Reference

- `addons` - List of Cluster Addons objects
    * `id` - (String)
    * `name` - (String) Addon name
    * `description` - (String) Addon description
    * `releases` - List of addon versions 
        + `id` - (String)
        + `version` - (String) Addon version






  






