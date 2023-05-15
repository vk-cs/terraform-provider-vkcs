resource "vkcs_vpnaas_endpoint_group" "group_1" {
	name = "Group 1"
	type = "cidr"
	endpoints = [
		"10.2.0.0/24",
		"10.3.0.0/24",
	]
}
