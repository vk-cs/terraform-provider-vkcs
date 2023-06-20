# Terraform Resource Schema Migrator

Migrates a Plugin SDK v2 resource schema to the identical Plugin Framework schema.

This tool

* Introspects a Plugin SDK v2 resource schema
* Generates Go code for the identical schema targeting the [Terraform Plugin Framework](https://github.com/hashicorp/terraform-plugin-framework)

Run `tfsdk2fw --help` to see all options.

### Example usage
```
go run ./helpers/tfsdk2fw -resource vkcs_kubernetes_node_group kubernetes NodeGroup vkcs/kubernetes/vkcs_kubernetes_node_group_fw.go
```