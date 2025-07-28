---
subcategory: "Data Platform"
layout: "vkcs"
page_title: "vkcs: vkcs_dataplatform_product"
description: |-
  Get information on VKCS Data Platform product.
---

# vkcs_dataplatform_product



## Example Usage

```terraform
data "vkcs_dataplatform_product" "spark" {
  product_name = "spark"
}
```

## Argument Reference
- `product_name` **required** *string* &rarr;  Product name

- `product_version` optional *string* &rarr;  Product version

- `region` optional *string* &rarr;  The region in which to obtain the Data platform client. If omitted, the `region` argument of the provider is used. Changing this creates a new resource.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `configs` 
  - `connections`  *list* &rarr;  Product connections configuration info
    - `is_required` *boolean* &rarr;  Is connection required

    - `plug` *string* &rarr;  Connection type

    - `position` *number* &rarr;  Connection position

    - `required_group` *string* &rarr;  Connection required group

    - `settings`  *list* &rarr;  Connection settings
      - `alias` *string* &rarr;  Setting alias

      - `default_value` *string* &rarr;  Setting default value

      - `is_require` *boolean* &rarr;  Is setting required

      - `is_sensitive` *boolean* &rarr;  Is setting sensitive

      - `regexp` *string* &rarr;  Setting validation regexp

      - `string_variation` *string* &rarr;  Available setting values



  - `settings`  *list* &rarr;  Product settings
    - `alias` *string* &rarr;  Setting alias

    - `default_value` *string* &rarr;  Setting default value

    - `is_require` *boolean* &rarr;  Is setting required

    - `is_sensitive` *boolean* &rarr;  Is setting sensitive

    - `regexp` *string* &rarr;  Setting validation regexp

    - `string_variation` *string* &rarr;  Available setting values




