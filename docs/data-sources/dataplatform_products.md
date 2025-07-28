---
subcategory: "Data Platform"
layout: "vkcs"
page_title: "vkcs: vkcs_dataplatform_products"
description: |-
  Get information on VKCS Data Platform products.
---

# vkcs_dataplatform_products



## Example Usage

```terraform
data "vkcs_dataplatform_products" "products" {}
```

## Argument Reference
- `region` optional *string* &rarr;  The region in which to obtain the Data platform client. If omitted, the `region` argument of the provider is used. Changing this creates a new resource.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `products`  *list* &rarr;  List of products information
  - `product_name` *string* &rarr;  Product name

  - `product_version` *string* &rarr;  Product version



