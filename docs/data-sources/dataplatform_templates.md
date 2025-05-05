---
subcategory: "Data Platform"
layout: "vkcs"
page_title: "vkcs: vkcs_dataplatform_templates"
description: |-
  Get information on a dataplatform templates.
---

# vkcs_dataplatform_templates



## Example Usage

## Argument Reference

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `cluster_templates`  *list* &rarr;  Cluster templates info.
  - `availability_zone` *string*

  - `configs`  &rarr;  Cluster template configs.
    - `common`  &rarr;  Common configs.
      - `maintenance`  &rarr;  Maintenance settings.
        - `backup`  &rarr;  Backup settings.
          - `differential`  &rarr;  Differential backup settings.
            - `backup_name_prefix` *string* &rarr;  Backup name prefix.

            - `backup_s3_bucket_name` *string* &rarr;  Backup S3 bucket name.

            - `creation_timeout` *number* &rarr;  Backup creation timeout.

            - `enabled` *boolean* &rarr;  Whether differential backup is enabled.

            - `keep_count` *number*

            - `keep_time` *number*

            - `start` *string* &rarr;  Backup schedule.


          - `full`  &rarr;  Full backup settings.
            - `backup_name_prefix` *string* &rarr;  Backup name prefix.

            - `backup_s3_bucket_name` *string* &rarr;  Backup S3 bucket name.

            - `creation_timeout` *number* &rarr;  Backup creation timeout.

            - `enabled` *boolean* &rarr;  Whether full backup is enabled.

            - `keep_count` *number*

            - `keep_time` *number*

            - `start` *string* &rarr;  Backup schedule.


          - `incremental`  &rarr;  Incremental backup settings.
            - `backup_name_prefix` *string* &rarr;  Backup name prefix.

            - `backup_s3_bucket_name` *string* &rarr;  Backup S3 bucket name.

            - `creation_timeout` *number* &rarr;  Backup creation timeout.

            - `enabled` *boolean* &rarr;  Whether incremental backup is enabled.

            - `keep_count` *number*

            - `keep_time` *number*

            - `start` *string* &rarr;  Backup schedule.



        - `duration` *number* &rarr;  Maintenance duration.

        - `start` *string* &rarr;  Maintenance cron schedule.




  - `created_at` *string* &rarr;  Cluster template creation timestamp.

  - `description` *string* &rarr;  Cluster template description.

  - `id` *string* &rarr;  Cluster template id.

  - `is_hidden` *boolean*

  - `multiaz` *boolean* &rarr;  Is multiple available zones mode enabled.

  - `name` *string* &rarr;  Cluster template name.

  - `pod_groups`  *list* &rarr;  Cluster pod groups.
    - `alias` *string* &rarr;  Alias.

    - `backref` *string* &rarr;  Backref.

    - `cluster_template_id` *string* &rarr;  Cluster template id.

    - `count` *number* &rarr;  Pod count.

    - `created_at` *string* &rarr;  Pod group creation timestamp.

    - `description` *string* &rarr;  Pod group name.

    - `id` *string* &rarr;  Pod group id.

    - `name` *string* &rarr;  Pod group name.

    - `node_processes` *string* &rarr;  Node processes.

    - `resource`  &rarr;  Resource settings.
      - `cpu_margin` *number* &rarr;  CPU margin settings.

      - `cpu_request` *string* &rarr;  CPU request settings.

      - `ram_margin` *number* &rarr;  RAM margin settings.

      - `ram_request` *string* &rarr;  RAM request settings.


    - `template_type` *string* &rarr;  Template type.

    - `volumes`  *map* &rarr;  Volumes settings.
      - `count` *number* &rarr;  Volume count.

      - `storage` *string* &rarr;  Storage size.

      - `storage_class_name` *string* &rarr;  Storage class name.



  - `presets`  *list* &rarr;  Presets info.
    - `name` *string* &rarr;  Preset name.

    - `pod_groups`  *list* &rarr;  Preset pod groups.
      - `count` *number* &rarr;  Pod count.

      - `meta` 
        - `create` 
          - `property` 
            - `multiplicator` 
              - `allow` *boolean*

              - `allow_add` *boolean*

              - `count` *number*

              - `is_read_only` *boolean*

              - `max` *number*

              - `min` *number*





      - `name` *string* &rarr;  Pod group name.

      - `resource`  &rarr;  Resource settings.
        - `cpu_request` *string* &rarr;  CPU request settings.

        - `ram_request` *string* &rarr;  RAM request settings.


      - `volumes`  *map* &rarr;  Volumes settings.
        - `count` *number* &rarr;  Volume count.

        - `storage` *string* &rarr;  Storage size.

        - `storage_class_name` *string* &rarr;  Storage class name.




  - `product_name` *string* &rarr;  Product name.

  - `product_type` *string* &rarr;  Product type.

  - `product_version` *string* &rarr;  Product version.

  - `template_type` *string* &rarr;  Template type.



