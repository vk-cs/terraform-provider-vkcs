output "clickhouse_settings_all" {
  value = {
    for s in data.vkcs_dataplatform_product.clickhouse.configs.settings :
    s.alias => s.default_value
    # Settings which first part of name ends with "_" are not available for tuning
    if !endswith(split(".", s.alias)[0], "_")
  }
}
