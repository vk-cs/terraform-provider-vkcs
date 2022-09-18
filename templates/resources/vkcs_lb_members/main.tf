resource "vkcs_lb_members" "members_1" {
	pool_id = "935685fb-a896-40f9-9ff4-ae531a3a00fe"

	member {
		address       = "192.168.199.23"
		protocol_port = 8080
	}

	member {
		address       = "192.168.199.24"
		protocol_port = 8080
	}
}
