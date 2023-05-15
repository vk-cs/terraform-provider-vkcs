resource "vkcs_lb_listener" "listener_1" {
	loadbalancer_id = "d9415786-5f1a-428b-b35f-2f1523e146d2"
	protocol        = "HTTP"
	protocol_port   = 8080

	insert_headers = {
		X-Forwarded-For = "true"
	}
}
