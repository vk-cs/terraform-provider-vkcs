output "clickhouse_pod_groups" {
  value = data.vkcs_dataplatform_template.clickhouse.pod_groups[*].name
}
