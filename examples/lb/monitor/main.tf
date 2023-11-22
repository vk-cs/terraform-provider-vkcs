resource "vkcs_lb_monitor" "worker_ping_life_checker" {
	name        = "worker-ping-life-checker-tf-example"
	pool_id     = vkcs_lb_pool.http.id
	type        = "PING"
	delay       = 20
	timeout     = 10
	max_retries = 5
}
