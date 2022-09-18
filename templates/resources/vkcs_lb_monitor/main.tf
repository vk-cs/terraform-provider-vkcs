resource "vkcs_lb_monitor" "monitor_1" {
	pool_id     = "${vkcs_lb_pool.pool_1.id}"
	type        = "PING"
	delay       = 20
	timeout     = 10
	max_retries = 5
}
