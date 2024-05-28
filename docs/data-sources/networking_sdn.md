---
subcategory: "Network"
layout: "vkcs"
page_title: "vkcs: vkcs_networking_sdn"
description: |-
  Get information on a VKCS SDNs.
---

# vkcs_networking_sdn

Use this data source to get a list of available VKCS SDNs in the current project. The first SDN is default. You do not have to specify default sdn argument in resources and datasources. You may specify non default SDN only for root resources such as `vkcs_networking_router`, `vkcs_networking_network`, `vkcs_networking_secgroup` (they do not depend on any other resource/datasource with sdn argument).

## Example Usage

```terraform
data "vkcs_networking_sdn" "sdn" {
}
```

## Argument Reference

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `id` *string* &rarr;  ID of the resource.

- `sdn` *string* &rarr;  Names of available VKCS SDNs in the current project.


