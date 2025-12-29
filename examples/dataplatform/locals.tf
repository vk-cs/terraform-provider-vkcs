locals {
  # Parse ClickHouse endpoints to find tcp one.
  clickhouse_tcp = one([
    for srv in vkcs_dataplatform_cluster.clickhouse.info.services :
    {
      host = regex(".*@([^:/]+):([0-9]+).*", srv.connection_string)[0]
      port = regex(".*@([^:/]+):([0-9]+).*", srv.connection_string)[1]
    }
    if srv.type == "connection_string" && endswith(srv.description, "tcp")
  ])
}
