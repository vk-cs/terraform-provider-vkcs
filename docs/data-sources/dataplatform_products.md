---
subcategory: "Data Platform"
layout: "vkcs"
page_title: "vkcs: vkcs_dataplatform_products"
description: |-
  Get information on a dataplatform products.
---

# vkcs_dataplatform_products



## Example Usage


## Argument Reference
- `products`  *set* &rarr;  Products info.
  - `configs` optional &rarr;  Product configuration.
    - `connections`  *set* &rarr;  Connections settings.
      - `settings`  *set* &rarr;  Additional connection settings.
        - `alias` read-only *string* &rarr;  Setting alias.

        - `backref` read-only *string* &rarr;  Setting backref.

        - `default_value` read-only *string* &rarr;  Setting default value.

        - `dependencies`  *set* &rarr;  Setting dependencies.
          - `anchors` read-only *set of* *string* &rarr;  Dependency anchors.

          - `kind` read-only *string* &rarr;  Dependency kind.


        - `is_require` read-only *boolean* &rarr;  Is setting required.

        - `is_sensitive` read-only *string* &rarr;  Is setting sensitive.

        - `max_value` read-only *string* &rarr;  Setting max value.

        - `min_value` read-only *string* &rarr;  Setting min value.

        - `policy`  &rarr;  Setting policy.
          - `add_on_create` read-only *boolean* &rarr;  Add on create policy.

          - `live_update` read-only *boolean* &rarr;  Live update policy.


        - `regexp` read-only *string* &rarr;  Setting regexp.

        - `string_variation` read-only *set of* *string* &rarr;  Setting string variations.

        - `validation` read-only *string* &rarr;  Setting validation.


      - `desc_i18n_key` read-only *string* &rarr;  Connection description i18n key.

      - `is_require` read-only *boolean* &rarr;  Is connection required.

      - `name_i18n_key` read-only *string* &rarr;  Connection name i18n key.

      - `plug` read-only *string* &rarr;  Connection plug.

      - `position` read-only *number* &rarr;  Connection position.

      - `required_group` read-only *string* &rarr;  Connection required group.


    - `cron_tabs`  *set* &rarr;  Cron tabs settings.
      - `settings`  *set* &rarr;  Additional cron settings.
        - `alias` read-only *string* &rarr;  Setting alias.

        - `backref` read-only *string* &rarr;  Setting backref.

        - `default_value` read-only *string* &rarr;  Setting default value.

        - `dependencies`  *set* &rarr;  Setting dependencies.
          - `anchors` read-only *set of* *string* &rarr;  Dependency anchors.

          - `kind` read-only *string* &rarr;  Dependency kind.


        - `is_require` read-only *boolean* &rarr;  Is setting required.

        - `is_sensitive` read-only *string* &rarr;  Is setting sensitive.

        - `max_value` read-only *string* &rarr;  Setting max value.

        - `min_value` read-only *string* &rarr;  Setting min value.

        - `policy`  &rarr;  Setting policy.
          - `add_on_create` read-only *boolean* &rarr;  Add on create policy.

          - `live_update` read-only *boolean* &rarr;  Live update policy.


        - `regexp` read-only *string* &rarr;  Setting regexp.

        - `string_variation` read-only *set of* *string* &rarr;  Setting string variations.

        - `validation` read-only *string* &rarr;  Setting validation.


      - `name` read-only *string* &rarr;  Cron tab name.

      - `required` read-only *boolean* &rarr;  Is cron required.

      - `start` read-only *string* &rarr;  Cron schedule.


    - `extensions`  *set* &rarr;  Extensions settings.
      - `settings`  *set* &rarr;  Additional extension settings.
        - `alias` read-only *string* &rarr;  Setting alias.

        - `backref` read-only *string* &rarr;  Setting backref.

        - `default_value` read-only *string* &rarr;  Setting default value.

        - `dependencies`  *set* &rarr;  Setting dependencies.
          - `anchors` read-only *set of* *string* &rarr;  Dependency anchors.

          - `kind` read-only *string* &rarr;  Dependency kind.


        - `is_require` read-only *boolean* &rarr;  Is setting required.

        - `is_sensitive` read-only *string* &rarr;  Is setting sensitive.

        - `max_value` read-only *string* &rarr;  Setting max value.

        - `min_value` read-only *string* &rarr;  Setting min value.

        - `policy`  &rarr;  Setting policy.
          - `add_on_create` read-only *boolean* &rarr;  Add on create policy.

          - `live_update` read-only *boolean* &rarr;  Live update policy.


        - `regexp` read-only *string* &rarr;  Setting regexp.

        - `string_variation` read-only *set of* *string* &rarr;  Setting string variations.

        - `validation` read-only *string* &rarr;  Setting validation.


      - `access_control` read-only *string* &rarr;  Extension access control mode.

      - `desc_i18n_key` read-only *string* &rarr;  Extension description i18n key.

      - `name_i18n_key` read-only *string* &rarr;  Extension name i18n key.

      - `persistent` read-only *boolean* &rarr;  Is extensions persistent.

      - `type` read-only *string* &rarr;  Extension type.

      - `version` read-only *string* &rarr;  Extension version.


    - `settings`  *set* &rarr;  Additional settings.
      - `alias` read-only *string* &rarr;  Setting alias.

      - `backref` read-only *string* &rarr;  Setting backref.

      - `default_value` read-only *string* &rarr;  Setting default value.

      - `dependencies`  *set* &rarr;  Setting dependencies.
        - `anchors` read-only *set of* *string* &rarr;  Dependency anchors.

        - `kind` read-only *string* &rarr;  Dependency kind.


      - `is_require` read-only *boolean* &rarr;  Is setting required.

      - `is_sensitive` read-only *string* &rarr;  Is setting sensitive.

      - `max_value` read-only *string* &rarr;  Setting max value.

      - `min_value` read-only *string* &rarr;  Setting min value.

      - `policy`  &rarr;  Setting policy.
        - `add_on_create` read-only *boolean* &rarr;  Add on create policy.

        - `live_update` read-only *boolean* &rarr;  Live update policy.


      - `regexp` read-only *string* &rarr;  Setting regexp.

      - `string_variation` read-only *set of* *string* &rarr;  Setting string variations.

      - `validation` read-only *string* &rarr;  Setting validation.


    - `user_accesses`  *set* &rarr;  User accesses settings.
      - `settings`  *set* &rarr;  Additional user access settings.
        - `alias` read-only *string* &rarr;  Setting alias.

        - `backref` read-only *string* &rarr;  Setting backref.

        - `default_value` read-only *string* &rarr;  Setting default value.

        - `dependencies`  *set* &rarr;  Setting dependencies.
          - `anchors` read-only *set of* *string* &rarr;  Dependency anchors.

          - `kind` read-only *string* &rarr;  Dependency kind.


        - `is_require` read-only *boolean* &rarr;  Is setting required.

        - `is_sensitive` read-only *string* &rarr;  Is setting sensitive.

        - `max_value` read-only *string* &rarr;  Setting max value.

        - `min_value` read-only *string* &rarr;  Setting min value.

        - `policy`  &rarr;  Setting policy.
          - `add_on_create` read-only *boolean* &rarr;  Add on create policy.

          - `live_update` read-only *boolean* &rarr;  Live update policy.


        - `regexp` read-only *string* &rarr;  Setting regexp.

        - `string_variation` read-only *set of* *string* &rarr;  Setting string variations.

        - `validation` read-only *string* &rarr;  Setting validation.




  - `product_name` read-only *string* &rarr;  Product name.

  - `product_version` read-only *string* &rarr;  Product name.


- `region` optional *string* &rarr;  The region in which to obtain the Data Platform client. If omitted, the `region` argument of the provider is used.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `id` *string* &rarr;  ID of the data source.


