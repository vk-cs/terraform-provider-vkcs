---
subcategory: "Data Platform"
layout: "vkcs"
page_title: "vkcs: vkcs_dataplatform_template"
description: |-
  Get information on VKCS Data Platform template.
---

# vkcs_dataplatform_template



## Example Usage

```terraform
data "vkcs_dataplatform_template" "spark" {
  product_name    = "spark"
  product_version = "3.5.1"
}
```

## Argument Reference
- `product_name` **required** *string* &rarr;  Product name.

- `product_version` optional *string* &rarr;  Product version.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `id` *string* &rarr;  ID of the cluster template.

- `name` *string* &rarr;  Name of the cluster template.

- `pod_groups`  *list* &rarr;  List of pod groups in the template.
    - `count` *number* &rarr;  Number of pods in the pod group.

    - `name` *string* &rarr;  Pod group name.

    - `resource`  &rarr;  Resource settings for the pod group.
        - `cpu_margin` *number* &rarr;  CPU margin for the pod group.

        - `cpu_request` *string* &rarr;  CPU request for the pod group.

        - `ram_margin` *number* &rarr;  RAM margin for the pod group.

        - `ram_request` *string* &rarr;  RAM request for the pod group.


    - `volumes`  *map* &rarr;  Volumes configuration for the pod group.
        - `count` *number* &rarr;  Volume count.

        - `storage` *string* &rarr;  Volume storage size.

        - `storage_class_name` *string* &rarr;  Storage class name for the volume.




