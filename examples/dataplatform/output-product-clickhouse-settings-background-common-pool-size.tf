output "clickhouse_settings_background_common_pool_size" {
  value = one([
    for s in data.vkcs_dataplatform_product.clickhouse.configs.settings :
    s if s.alias == "clickhouse.background_common_pool_size"
  ])
}
