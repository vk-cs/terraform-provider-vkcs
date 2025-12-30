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

# ClickHouse
```terraform
resource "vkcs_dataplatform_cluster" "clickhouse" {
  name            = "clickhouse-tf-guide"
  description     = "ClickHouse example instance from Data Platform guide."
  product_name    = "clickhouse"
  product_version = "25.3.0"

  network_id        = vkcs_networking_network.db.id
  subnet_id         = vkcs_networking_subnet.db.id
  availability_zone = "GZ1"
  # Enable public access to simplify testing of the product.
  floating_ip_pool = "auto"

  pod_groups = [
    # Omit settings for clickhouseKeeper pod group to illustrate
    # how Data Platform handle this with default settings.
    # NOTE: If you omit settings for a pod group you cannot scale
    # the pod group later.
    # Increase ram_request and storage values for clickhouse pod group
    # against default settings.
    {
      count = 3
      name  = "clickhouse"
      resource = {
        cpu_request = "2.0"
        ram_request = "8.0"
      }
      volumes = {
        data = {
          storage_class_name = "ceph-ssd"
          storage            = "150"
          count              = 1
        }
      }
    },
  ]
  configs = {
    settings = [
      # Increase value of the setting against default one.
      {
        alias = "clickhouse.background_common_pool_size"
        value = 10
      },
    ]
    users = [
      {
        username = "owner"
        password = random_password.clickhouse_owner.result
        role     = "dbOwner"
      },
      {
        username = "trino"
        password = random_password.clickhouse_trino.result
        role     = "readOnly"
      },
    ]
    warehouses = [
      # Define database name.
      {
        name = "clickhouse"
      },
    ]
    maintenance = {
      # Set start om maintenance the same as start of full backup.
      # Otherwise you get unpredictable behavior of interaction between
      # Terraform, VKCS Terraform provider and Data Platform API.
      start = "0 1 * * 0"
      backup = {
        full = {
          keep_count = 5
          start      = "0 1 * * 0"
        }
        incremental = {
          keep_count = 7
          start      = "0 1 * * 1-6"
        }
      }
    }
  }

  # If you create networking in the same bundle of resources with Data Platform resource
  # add dependency on corresponding vkcs_networking_router_interface resource.
  # However this is not required if you set up networking witth terraform-vkcs-network module.
  depends_on = [vkcs_networking_router_interface.db]
}
```

See more examples on [GitHub](https://github.com/vk-cs/terraform-provider-vkcs/tree/master/examples/dataplatform/cluster).

Refer to the `Setting up Data Platform products` guide for details of using the resource in combination with Data Platform datasources.

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

- `product_name` **required** *string* &rarr;  Name of the product. Changing this creates a new resource.

- `product_version` **required** *string* &rarr;  Version of the product. Changing this creates a new resource.

- `availability_zone` optional *string* &rarr;  Availability zone to create cluster in. Changing this creates a new resource.

- `cluster_template_id` optional *string* &rarr;  ID of the cluster template. Changing this creates a new resource.

- `description` optional *string* &rarr;  Cluster description.

- `floating_ip_pool` optional *string* &rarr;  Floating IP pool ID. Use `auto` for autoselect. Changing this creates a new resource.

- `multiaz` optional *boolean* &rarr;  Enables multi az support. Changing this creates a new resource.

- `pod_groups`  *list* &rarr;  Cluster pod groups. Changing this creates a new resource.
    - `name` **required** *string* &rarr;  Pod group name. Changing this creates a new resource.

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