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



    - `crontabs`  *list* &rarr;  Product crontabs
        - `name` *string* &rarr;  Crontab name

        - `required` *boolean* &rarr;  Crontab required

        - `settings`  *list* &rarr;  Crontabs settings
            - `alias` *string* &rarr;  Setting alias

            - `default_value` *string* &rarr;  Setting default value

            - `is_require` *boolean* &rarr;  Is setting required

            - `is_sensitive` *boolean* &rarr;  Is setting sensitive

            - `regexp` *string* &rarr;  Setting validation regexp

            - `string_variation` *string* &rarr;  Available setting values


        - `start` *string* &rarr;  Crontab start


    - `settings`  *list* &rarr;  Product settings
        - `alias` *string* &rarr;  Setting alias

        - `default_value` *string* &rarr;  Setting default value

        - `is_require` *boolean* &rarr;  Is setting required

        - `is_sensitive` *boolean* &rarr;  Is setting sensitive

        - `regexp` *string* &rarr;  Setting validation regexp

        - `string_variation` *string* &rarr;  Available setting values


    - `user_roles`  *list* &rarr;  User roles list
        - `name` *string* &rarr;  User role name




