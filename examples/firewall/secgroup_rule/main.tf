resource "vkcs_networking_secgroup_rule" "etcd_app_clients" {
  description       = "etcd app clients rule"
  security_group_id = vkcs_networking_secgroup.etcd.id
  direction         = "ingress"
  protocol          = "tcp"
  port_range_max    = 2379
  port_range_min    = 2379
  remote_ip_prefix  = vkcs_networking_subnet.app.cidr
}
