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
- `products`  *set* &rarr;  Cluster templates info.
  - `configs`  &rarr;  Cluster template configs.
    - `common`  &rarr;  Common configs.
      - `maintenance`  &rarr;  Maintenance settings.
        - `backup`  &rarr;  Backup settings.
          - `differential`  &rarr;  Differential backup settings.
            - `backup_name_prefix` read-only *string* &rarr;  Backup name prefix.

            - `backup_s3_bucket_name` read-only *string* &rarr;  Backup S3 bucket name.

            - `creation_timeout` read-only *number* &rarr;  Backup creation timeout.

            - `enabled` read-only *boolean* &rarr;  Whether differential backup is enabled.

            - `keep_count` read-only *number* &rarr;  Backup keep count.

            - `keep_time` read-only *number* &rarr;  Backup keep time.

            - `start` read-only *string* &rarr;  Backup schedule.


          - `full`  &rarr;  Full backup settings.
            - `backup_name_prefix` read-only *string* &rarr;  Backup name prefix.

            - `backup_s3_bucket_name` read-only *string* &rarr;  Backup S3 bucket name.

            - `creation_timeout` read-only *number* &rarr;  Backup creation timeout.

            - `enabled` read-only *boolean* &rarr;  Whether full backup is enabled.

            - `keep_count` read-only *number* &rarr;  Backup keep count.

            - `keep_time` read-only *number* &rarr;  Backup keep time.

            - `start` read-only *string* &rarr;  Backup schedule.


          - `incremental`  &rarr;  Incremental backup settings.
            - `backup_name_prefix` read-only *string* &rarr;  Backup name prefix.

            - `backup_s3_bucket_name` read-only *string* &rarr;  Backup S3 bucket name.

            - `creation_timeout` read-only *number* &rarr;  Backup creation timeout.

            - `enabled` read-only *boolean* &rarr;  Whether incremental backup is enabled.

            - `keep_count` read-only *number* &rarr;  Backup keep count.

            - `keep_time` read-only *number* &rarr;  Backup keep time.

            - `start` read-only *string* &rarr;  Backup schedule.



        - `duration` read-only *number* &rarr;  Maintenance duration.

        - `start` read-only *string* &rarr;  Maintenance cron schedule.




  - `created_at` read-only *string* &rarr;  Cluster template creation timestamp.

  - `description` read-only *string* &rarr;  Cluster template name.

  - `id` read-only *string* &rarr;  Cluster template id.

  - `multiaz` read-only *boolean* &rarr;  Is multiple available zones mode enabled.

  - `name` read-only *string* &rarr;  Cluster template name.

  - `pod_groups`  *set* &rarr;  Cluster pod groups.
    - `alias` read-only *string* &rarr;  Alias.

    - `backref` read-only *string* &rarr;  Backref.

    - `cluster_template_id` read-only *string* &rarr;  Cluster template id.

    - `count` read-only *number* &rarr;  Pod count.

    - `created_at` read-only *string* &rarr;  Pod group creation timestamp.

    - `description` read-only *string* &rarr;  Pod group name.

    - `id` read-only *string* &rarr;  Pod group id.

    - `name` read-only *string* &rarr;  Pod group name.

    - `node_processes` read-only *set of* *string* &rarr;  Node processes.

    - `resource`  &rarr;  Resource settings.
      - `cpu_margin` read-only *number* &rarr;  CPU margin settings.

      - `cpu_request` read-only *string* &rarr;  CPU request settings.

      - `ram_margin` read-only *number* &rarr;  RAM margin settings.

      - `ram_request` read-only *string* &rarr;  RAM request settings.


    - `template_type` read-only *string* &rarr;  Template type.

    - `volumes`  *set* &rarr;  Volumes settings.
      - `count` read-only *number* &rarr;  Volume count.

      - `storage` read-only *string* &rarr;  Storage size.

      - `storage_class_name` read-only *string* &rarr;  Storage class name.



  - `presets`  *set* &rarr;  Presets info.
    - `name` read-only *string* &rarr;  Preset name.

    - `pod_groups`  *set* &rarr;  Preset pod groups.
      - `count` read-only *number* &rarr;  Pod count.

      - `name` read-only *string* &rarr;  Pod group name.

      - `resource`  &rarr;  Resource settings.
        - `cpu_request` read-only *string* &rarr;  CPU request settings.

        - `ram_request` read-only *string* &rarr;  RAM request settings.


      - `volumes`  *set* &rarr;  Volumes settings.
        - `count` read-only *number* &rarr;  Volume count.

        - `storage` read-only *string* &rarr;  Storage size.

        - `storage_class_name` read-only *string* &rarr;  Storage class name.




  - `product_name` read-only *string* &rarr;  Product name.

  - `product_type` read-only *string* &rarr;  Product type.

  - `product_version` read-only *string* &rarr;  Product version.

  - `template_type` read-only *string* &rarr;  Template type.


- `region` optional *string* &rarr;  The region in which to obtain the Data Platform client. If omitted, the `region` argument of the provider is used.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `id` *string* &rarr;  ID of the data source.


