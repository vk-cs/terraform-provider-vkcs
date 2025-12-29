output "trino_connection_clickhouse" {
  value = one([
    for c in data.vkcs_dataplatform_product.trino.configs.connections :
    c.settings
    if c.plug == "clickhouse"
  ])
}
