# Create networks
resource "vkcs_networking_network" "app" {
  name        = "app-tf-example"
  description = "Application network"
}

resource "vkcs_networking_subnet" "app" {
  name       = "app-tf-example"
  network_id = vkcs_networking_network.app.id
  cidr       = "192.168.199.0/24"
}

resource "vkcs_networking_network" "db" {
  name        = "db-tf-example"
  description = "Database network"
}

resource "vkcs_networking_subnet" "db" {
  name       = "db-tf-example"
  network_id = vkcs_networking_network.db.id
  cidr       = "192.168.166.0/24"
}

# Get external network with Inernet access
data "vkcs_networking_network" "extnet" {
  name = "ext-net"
}

# Create a router to connect netwoks
resource "vkcs_networking_router" "router" {
  name = "router-tf-example"
  # Connect router to Internet
  external_network_id = data.vkcs_networking_network.extnet.id
}

# Connect networks to the router
resource "vkcs_networking_router_interface" "app" {
  router_id = vkcs_networking_router.router.id
  subnet_id = vkcs_networking_subnet.app.id
}

resource "vkcs_networking_router_interface" "db" {
  router_id = vkcs_networking_router.router.id
  subnet_id = vkcs_networking_subnet.db.id
}

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
