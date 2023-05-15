resource "vkcs_networking_network" "network_1" {
	name           = "network_1"
	admin_state_up = "true"
}

resource "vkcs_networking_subnet" "subnet_1" {
	name       = "subnet_1"
	cidr       = "192.168.199.0/24"
	network_id = "${vkcs_networking_network.network_1.id}"
}

resource "vkcs_lb_loadbalancer" "loadbalancer_1" {
	name          = "loadbalancer_1"
	vip_subnet_id = "${vkcs_networking_subnet.subnet_1.id}"
}

resource "vkcs_lb_listener" "listener_1" {
	name            = "listener_1"
	protocol        = "HTTP"
	protocol_port   = 8080
	loadbalancer_id = "${vkcs_lb_loadbalancer.loadbalancer_1.id}"
}

resource "vkcs_lb_pool" "pool_1" {
	name            = "pool_1"
	protocol        = "HTTP"
	lb_method       = "ROUND_ROBIN"
	loadbalancer_id = "${vkcs_lb_loadbalancer.loadbalancer_1.id}"
}

resource "vkcs_lb_l7policy" "l7policy_1" {
	name             = "test"
	action           = "REDIRECT_TO_POOL"
	description      = "test l7 policy"
	position         = 1
	listener_id      = "${vkcs_lb_listener.listener_1.id}"
	redirect_pool_id = "${vkcs_lb_pool.pool_1.id}"
}
