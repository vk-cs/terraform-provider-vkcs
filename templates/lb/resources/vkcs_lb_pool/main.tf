resource "vkcs_lb_pool" "pool_1" {
	protocol    = "HTTP"
	lb_method   = "ROUND_ROBIN"
	listener_id = "d9415786-5f1a-428b-b35f-2f1523e146d2"

	persistence {
		type        = "APP_COOKIE"
		cookie_name = "testCookie"
	}
}
