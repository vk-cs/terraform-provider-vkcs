resource "vkcs_lb_member" "member_1" {
	address       = "192.168.199.23"
	pool_id       = "935685fb-a896-40f9-9ff4-ae531a3a00fe"
	protocol_port = 8080
}
