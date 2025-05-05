---
subcategory: "Data Platform"
layout: "vkcs"
page_title: "vkcs: vkcs_dataplatform_cluster"
description: |-
  Manages a dataplatform cluster
---

# vkcs_dataplatform_cluster



## Example Usage

# Spark
```terraform
data "vkcs_dataplatform_templates" "cluster-templates" {
}

resource "vkcs_dataplatform_cluster" "basic-opensearch" {
  cluster_template_id = data.vkcs_dataplatform_templates.cluster-templates.templates[index(data.vkcs_dataplatform_templates.cluster-templates.*.product_name, "spark")].id
  name                = "tf-basic-spark"
  network_id          = "d2fad739-1b10-4dc8-9b2c-c246d7a7cc69"
  subnet_id           = "3a744943-fcc1-4a85-a96b-3dc4fff71885"
  product_name        = "spark"
  product_version     = "3.5.1"

  availability_zone = "UD2"
  configs = {
    settings = [
      {
        alias = "sparkproxy.spark_version"
        value = "spark-py-3.5.1:v3.5.1.2"
      }
    ]
    users = [
      {
        user     = "user"
        password = "somepa55word!"
      }
    ]
    maintenance = {
      start = "0 0 1 * *"
    }
    warehouses = [
      {
        name = "db_customer"
        users = [
          "user"
        ]
        connections = [
          {
            name = "s3_int"
            plug = "s3-int"
            settings = [
              {
                alias = "s3_bucket"
                value = "really-cool-bucket"
              },
              {
                alias = "s3_folder"
                value = "folder"
              }
            ]
          },
          {
            name = "postgres"
            plug = "postgresql"
            settings = [
              {
                alias = "db_name"
                value = "db"
              },
              {
                alias = "hostname"
                value = "database.com:5432"
              },
              {
                alias = "username"
                value = "db"
              },
              {
                alias = "password"
                value = "db"
              }
            ]
          }
        ]
      }
    ]
  }
  pod_groups = [
    {
      count = 1
      resource = {
        cpu_request = "8"
        ram_request = "8"
      }
      pod_group_template_id = "6a8e3515-d0d6-40f9-826e-e33dbe141485"
    },
    {
      count = 1
      resource = {
        cpu_request = "0.5"
        ram_request = "0.5"
      }
      volumes = [
        {
          type = "data"
          storage_class_name = "ceph"
          storage = "5"
          count = 1
        }
      ]
      pod_group_template_id = "498c4bf6-1e3e-4a06-b8b7-d60337e85dc1"
    }
  ]
}
```

## Argument Reference
- `cluster_template_id` **required** *string* &rarr;  ID of the cluster template.

- `configs` ***required*** &rarr;  Product configuration.
  - `maintenance` ***required*** &rarr;  Maintenance settings.
    - `backup` optional &rarr;  Backup settings.
      - `differential` optional &rarr;  Differential backup settings.
        - `start` **required** *string* &rarr;  Differential backup schedule.

        - `keep_count` optional *number*

        - `keep_time` optional *number*

        - `enabled` read-only *boolean* &rarr;  Whether differential backup is enabled.


      - `full` optional &rarr;  Full backup settings.
        - `start` **required** *string* &rarr;  Full backup schedule.

        - `keep_count` optional *number*

        - `keep_time` optional *number*

        - `enabled` read-only *boolean* &rarr;  Whether full backup is enabled.


      - `incremental` optional &rarr;  Incremental backup settings.
        - `start` **required** *string* &rarr;  Incremental backup schedule.

        - `keep_count` optional *number*

        - `keep_time` optional *number*

        - `enabled` read-only *boolean* &rarr;  Whether incremental backup is enabled.



    - `crontabs`  *list* &rarr;  Cron tabs settings.
      - `name` **required** *string* &rarr;  Cron tab name.

      - `settings`  *list* &rarr;  Additional cron settings.
        - `alias` **required** *string* &rarr;  Setting alias.

        - `value` **required** *string* &rarr;  Setting value.


      - `start` optional *string* &rarr;  Cron tab schedule.

      - `required` read-only *boolean* &rarr;  Whether cron tab is required.


    - `start` optional *string* &rarr;  Maintenance cron schedule.


  - `features` optional &rarr;  Product features.
    - `volume_autoresize` optional &rarr;  Volume autoresize options.
      - `data` optional &rarr;  Data volume options.
        - `enabled` optional *boolean* &rarr;  Enables option.

        - `max_scale_size` optional *number* &rarr;  Maximum scale size.

        - `scale_step_size` optional *number* &rarr;  Scale step size.

        - `size_scale_threshold` optional *number* &rarr;  Size scale threshold.


      - `wal` optional &rarr;  Data volume options.
        - `enabled` optional *boolean* &rarr;  Enables option.

        - `max_scale_size` optional *number* &rarr;  Maximum scale size.

        - `scale_step_size` optional *number* &rarr;  Scale step size.

        - `size_scale_threshold` optional *number* &rarr;  Size scale threshold.




  - `settings`  *list* &rarr;  Additional common settings.
    - `alias` **required** *string* &rarr;  Setting alias.

    - `value` **required** *string* &rarr;  Setting value.


  - `users`  *list* &rarr;  Users settings.
    - `password` **required** *string* &rarr;  Password.

    - `username` **required** *string* &rarr;  Username.

    - `access` optional &rarr;  Access settings.
      - `settings`  *list* &rarr;  Access users settings.
        - `alias` **required** *string* &rarr;  Setting alias.

        - `value` **required** *string* &rarr;  Setting value.


      - `id` read-only *string* &rarr;  Access ID.


    - `role` optional *string* &rarr;  User role.

    - `created_at` read-only *string* &rarr;  User creation timestamp.

    - `id` read-only *string* &rarr;  User ID.


  - `warehouses`  *list* &rarr;  Warehouses settings.
    - `connections`  *list* &rarr;  Warehouse connections.
      - `name` **required** *string* &rarr;  Connection name.

      - `plug` **required** *string* &rarr;  Connection plug.

      - `settings`  *list* &rarr;  Additional warehouse settings.
        - `alias` **required** *string* &rarr;  Setting alias.

        - `value` **required** *string* &rarr;  Setting value.


      - `created_at` read-only *string* &rarr;  Connection creation timestamp.

      - `id` read-only *string* &rarr;  Connection ID.


    - `extensions`  *list* &rarr;  Warehouse extensions.
      - `type` **required** *string* &rarr;  Extension type.

      - `settings`  *list* &rarr;  Additional extension settings.
        - `alias` **required** *string* &rarr;  Setting alias.

        - `value` **required** *string* &rarr;  Setting value.


      - `version` optional *string* &rarr;  Extension version.

      - `created_at` read-only *string* &rarr;  Extension creation timestamp.

      - `id` read-only *string* &rarr;  Extension ID


    - `name` optional *string* &rarr;  Warehouse name.

    - `users` optional *string* &rarr;  Warehouse users.

    - `id` read-only *string* &rarr;  Warehouse ID.



- `name` **required** *string* &rarr;  Name of the cluster.

- `network_id` **required** *string* &rarr;  ID of the cluster network.

- `product_name` **required** *string* &rarr;  Name of the product.

- `product_version` **required** *string* &rarr;  Version of the product.

- `availability_zone` optional *string* &rarr;  Availability zone to create cluster in.

- `cluster_id` optional *string*

- `description` optional *string* &rarr;  Cluster description.

- `floating_ip_pool` optional *string* &rarr;  Floating IP pool ID.

- `multiaz` optional *boolean* &rarr;  Enables multi az support.

- `pod_groups`  *list* &rarr;  Cluster pod groups.
  - `pod_group_template_id` **required** *string* &rarr;  Pod group template ID.

  - `count` optional *number* &rarr;  Pod count.

  - `floating_ip_pool` optional *string* &rarr;  Floating IP pool ID.

  - `node_processes` optional *string* &rarr;  Node processes.

  - `resource` optional &rarr;  Resource request settings.
    - `cpu_request` optional *string* &rarr;  Resource request settings.

    - `ram_request` optional *string* &rarr;  RAM request settings.

    - `cpu_limit` read-only *string* &rarr;  CPU limit.

    - `ram_limit` read-only *string* &rarr;  RAM limit settings.


  - `volumes`  *map* &rarr;  Volumes settings.
    - `count` **required** *number* &rarr;  Volume count.

    - `storage` **required** *string* &rarr;  Storage size.

    - `storage_class_name` **required** *string* &rarr;  Storage class name.


  - `alias` read-only *string* &rarr;  Pod group alias.

  - `availability_zone` read-only *string*

  - `id` read-only *string* &rarr;  Pod group ID.

  - `name` read-only *string* &rarr;  Pod group name.


- `region` optional *string* &rarr;  The region in which to obtain the Data platform client. If omitted, the `region` argument of the provider is used. Changing this creates a new resource.

- `stack_id` optional *string* &rarr;  ID of the cluster stack.

- `subnet_id` optional *string* &rarr;  ID of the cluster subnet.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `created_at` *string* &rarr;  Cluster creation timestamp.

- `id` *string* &rarr;  ID of the cluster.

- `info`  &rarr;  Cluster info.
  - `services`  *list* &rarr;  Cluster services info.
    - `connection_string` *string* &rarr;  Service connection string.

    - `description` *string* &rarr;  Service description.

    - `exposed` *boolean* &rarr;  Is service exposed.

    - `type` *string* &rarr;  Service type.



- `product_type` *string* &rarr;  Type of the product.

- `project_id` *string*

- `status` *string* &rarr;  Cluster status.

- `status_description` *string* &rarr;  Cluster status description.


