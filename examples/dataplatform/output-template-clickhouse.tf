output "clickhouse_pod_group" {
  value = one([
    for g in data.vkcs_dataplatform_template.clickhouse.pod_groups :
    g if g.name == "clickhouse"
  ])
}
