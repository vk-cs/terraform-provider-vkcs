# Create security groups to define networking access
resource "vkcs_networking_secgroup" "admin" {
  name        = "admin-tf-example"
  description = "Admin access"
}

resource "vkcs_networking_secgroup_rule" "ssh" {
  description       = "SSH rule"
  security_group_id = vkcs_networking_secgroup.admin.id
  direction         = "ingress"
  protocol          = "tcp"
  # Specify SSH port
  port_range_max = 22
  port_range_min = 22
  # Allow access from any sources
  remote_ip_prefix = "0.0.0.0/0"
}

resource "vkcs_networking_secgroup_rule" "ping" {
  description       = "Ping rule"
  security_group_id = vkcs_networking_secgroup.admin.id
  direction         = "ingress"
  protocol          = "icmp"
  remote_ip_prefix  = "0.0.0.0/0"
}

resource "vkcs_networking_secgroup" "http" {
  name        = "http-tf-example"
  description = "HTTP access"
}

resource "vkcs_networking_secgroup_rule" "http" {
  description       = "HTTP rule"
  security_group_id = vkcs_networking_secgroup.http.id
  direction         = "ingress"
  protocol          = "tcp"
  port_range_max    = 80
  port_range_min    = 80
  # Allow access from application network only
  remote_ip_prefix = vkcs_networking_subnet.app.cidr
}

resource "vkcs_networking_secgroup_rule" "http_alter" {
  description       = "Alternative HTTP rule"
  security_group_id = vkcs_networking_secgroup.http.id
  direction         = "ingress"
  protocol          = "tcp"
  port_range_max    = 8080
  port_range_min    = 8080
  remote_ip_prefix  = vkcs_networking_subnet.app.cidr
}

resource "vkcs_networking_secgroup" "etcd" {
  name        = "etcd-tf-example"
  description = "etcd service"
}

resource "vkcs_networking_secgroup_rule" "etcd_app_clients" {
  description       = "etcd app clients rule"
  security_group_id = vkcs_networking_secgroup.etcd.id
  direction         = "ingress"
  protocol          = "tcp"
  port_range_max    = 2379
  port_range_min    = 2379
  remote_ip_prefix  = vkcs_networking_subnet.app.cidr
}

resource "vkcs_networking_secgroup_rule" "etcd_db_clients" {
  description       = "etcd db clients rule"
  security_group_id = vkcs_networking_secgroup.etcd.id
  direction         = "ingress"
  protocol          = "tcp"
  port_range_max    = 2379
  port_range_min    = 2379
  remote_ip_prefix  = vkcs_networking_subnet.db.cidr
}

resource "vkcs_networking_secgroup_rule" "etcd_peers" {
  description       = "etcd peer communication rule"
  security_group_id = vkcs_networking_secgroup.etcd.id
  direction         = "ingress"
  protocol          = "tcp"
  port_range_max    = 2380
  port_range_min    = 2380
  remote_group_id   = vkcs_networking_secgroup.etcd.id
}
