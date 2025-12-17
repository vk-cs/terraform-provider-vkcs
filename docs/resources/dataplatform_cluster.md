---
subcategory: "Data Platform"
layout: "vkcs"
page_title: "vkcs: vkcs_dataplatform_cluster"
description: |-
  Manages a dataplatform cluster
---

# vkcs_dataplatform_cluster

~> **Note:** Dataplatform cluster resource is currently in beta status.


## Example Usage

# Spark
```terraform
resource "vkcs_dataplatform_cluster" "basic_spark" {
  name            = "tf-basic-spark"
  description     = "tf-basic-description"
  network_id      = vkcs_networking_network.db.id
  subnet_id       = vkcs_networking_subnet.db.id
  product_name    = "spark"
  product_version = "3.5.1"

  availability_zone = "GZ1"
  configs = {
    settings = [
      {
        alias = "sparkproxy.spark_version"
        value = "spark-py-3.5.1:v3.5.1.2"
      }
    ]
    maintenance = {
      start = "0 0 1 * *"
    }
    warehouses = [
      {
        name = "spark"
        connections = [
          {
            name = "s3_int"
            plug = "s3-int"
            settings = [
              {
                alias = "s3_bucket"
                value = local.s3_bucket
              },
              {
                alias = "s3_folder"
                value = "tfexample-folder"
              }
            ]
          },
          {
            name = "postgres"
            plug = "postgresql"
            settings = [
              {
                alias = "db_name"
                value = vkcs_db_database.postgres_db.name
              },
              {
                alias = "hostname"
                value = "${vkcs_db_instance.db_instance.ip[0]}:5432"
              },
              {
                alias = "username"
                value = vkcs_db_user.postgres_user.name
              },
              {
                alias = "password"
                value = vkcs_db_user.postgres_user.password
              }
            ]
          }
        ]
      }
    ]
  }
  pod_groups = [
    {
      name  = "sparkconnect"
      count = 1
      resource = {
        cpu_request = "10"
        ram_request = "10"
      }
    },
    {
      name  = "sparkhistory"
      count = 1
      resource = {
        cpu_request = "0.5"
        ram_request = "1"
      }
      volumes = {
        "data" = {
          storage_class_name = "ceph-ssd"
          storage            = "5"
          count              = 1
        }
      }
    }
  ]

  depends_on = [vkcs_networking_router_interface.db]
}
```

## Argument Reference
- `configs` ***required*** &rarr;  Product configuration.
    - `maintenance` ***required*** &rarr;  Maintenance settings.
        - `backup` optional &rarr;  Backup settings.
            - `differential` optional &rarr;  Differential backup settings.
                - `start` **required** *string* &rarr;  Differential backup schedule. Defined in UTC.

                - `keep_count` optional *number*

                - `keep_time` optional *number*

                - `enabled` read-only *boolean* &rarr;  Whether differential backup is enabled.


            - `full` optional &rarr;  Full backup settings.
                - `start` **required** *string* &rarr;  Full backup schedule. Defined in UTC. <br>**Note:** `configs.maintenance.backup.full.start` must be equal to `configs.maintenance.start`.

                - `keep_count` optional *number*

                - `keep_time` optional *number*

                - `enabled` read-only *boolean* &rarr;  Whether full backup is enabled.


            - `incremental` optional &rarr;  Incremental backup settings.
                - `start` **required** *string* &rarr;  Incremental backup schedule. Defined in UTC.

                - `keep_count` optional *number*

                - `keep_time` optional *number*

                - `enabled` read-only *boolean* &rarr;  Whether incremental backup is enabled.



        - `crontabs`  *list* &rarr;  Cron tabs settings.
            - `name` **required** *string* &rarr;  Cron tab name.

            - `settings`  *list* &rarr;  Additional cron settings.
                - `alias` **required** *string* &rarr;  Setting alias.

                - `value` **required** *string* &rarr;  Setting value.


            - `start` optional *string* &rarr;  Cron tab schedule. Defined in UTC.

            - `id` read-only *string*

            - `required` read-only *boolean* &rarr;  Whether cron tab is required.


        - `start` optional *string* &rarr;  Maintenance cron schedule. Defined in UTC.


    - `settings`  *list* &rarr;  Additional common settings.
        - `alias` **required** *string* &rarr;  Setting alias.

        - `value` **required** *string* &rarr;  Setting value.


    - `users`  *list* &rarr;  Users settings.
        - `password` **required** sensitive *string* &rarr;  Password. Changing this creates a new resource.

        - `username` **required** *string* &rarr;  Username

        - `role` optional *string* &rarr;  User role. Changing this creates a new resource.

        - `created_at` read-only *string*

        - `id` read-only *string*


    - `warehouses`  *list* &rarr;  Warehouses settings. Changing this creates a new resource.
        - `connections`  *list* &rarr;  Warehouse connections.
            - `name` **required** *string* &rarr;  Connection name.

            - `plug` **required** *string* &rarr;  Connection plug.

            - `settings`  *list* &rarr;  Additional warehouse settings.
                - `alias` **required** *string* &rarr;  Setting alias.

                - `value` **required** *string* &rarr;  Setting value.


            - `created_at` read-only *string* &rarr;  Connection creation timestamp.

            - `id` read-only *string* &rarr;  Connection ID.


        - `name` optional *string* &rarr;  Warehouse name. Changing this creates a new resource.

        - `id` read-only *string* &rarr;  Warehouse ID.



- `name` **required** *string* &rarr;  Name of the cluster.

- `network_id` **required** *string* &rarr;  ID of the cluster network. Changing this creates a new resource.

- `product_name` **required** *string* &rarr;  Name of the product.

- `product_version` **required** *string* &rarr;  Version of the product. Changing this creates a new resource.

- `availability_zone` optional *string* &rarr;  Availability zone to create cluster in. Changing this creates a new resource.

- `cluster_template_id` optional *string* &rarr;  ID of the cluster template. Changing this creates a new resource.

- `description` optional *string* &rarr;  Cluster description.

- `floating_ip_pool` optional *string* &rarr;  Floating IP pool ID. Use `auto` for autoselect. Changing this creates a new resource.

- `multiaz` optional *boolean* &rarr;  Enables multi az support. Changing this creates a new resource.

- `pod_groups`  *list* &rarr;  Cluster pod groups. Changing this creates a new resource.
    - `name` **required** *string* &rarr;  Pod group name.

    - `count` optional *number* &rarr;  Pod count.

    - `floating_ip_pool` optional *string* &rarr;  Floating IP pool ID. Changing this creates a new resource.

    - `resource` optional &rarr;  Resource request settings. Changing this creates a new resource.
        - `cpu_request` optional *string* &rarr;  Resource request settings. Changing this creates a new resource.

        - `ram_request` optional *string* &rarr;  RAM request settings. Changing this creates a new resource.

        - `cpu_limit` read-only *string* &rarr;  CPU limit.

        - `ram_limit` read-only *string* &rarr;  RAM limit settings.


    - `volumes`  *map* &rarr;  Volumes settings. Changing this creates a new resource.
        - `count` **required** *number* &rarr;  Volume count. Changing this creates a new resource.

        - `storage` **required** *string* &rarr;  Storage size.

        - `storage_class_name` **required** *string* &rarr;  Storage class name. Changing this creates a new resource.


    - `alias` read-only *string* &rarr;  Pod group alias.

    - `availability_zone` read-only *string*

    - `id` read-only *string* &rarr;  Pod group ID.


- `region` optional *string* &rarr;  The region in which to obtain the Data platform client. If omitted, the `region` argument of the provider is used. Changing this creates a new resource.

- `stack_id` optional *string* &rarr;  ID of the cluster stack. Changing this creates a new resource.

- `subnet_id` optional *string* &rarr;  ID of the cluster subnet. Changing this creates a new resource.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `created_at` *string* &rarr;  Cluster creation timestamp.

- `id` *string* &rarr;  ID of the cluster.

- `info`  &rarr;  Application info
    - `services`  *list* &rarr;  Application services info
        - `connection_string` *string* &rarr;  Service connection string

        - `description` *string* &rarr;  Service description

        - `exposed` *boolean* &rarr;  Whether service is exposed

        - `type` *string* &rarr;  Service type



- `product_type` *string* &rarr;  Type of the product.



## Import

A Dataplaform cluster can be imported using the `id`, e.g.
```shell
terraform import vkcs_dataplatform_cluster.mycluster 83e08ade-c7cd-4382-8ee2-d297abbfc8d0
```

**Note:** Please, use `IMPORTED_PASSWORD` as value for external users' passwords.