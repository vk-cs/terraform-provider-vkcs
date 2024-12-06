data "vkcs_lb_loadbalancer" "app" {
  id = vkcs_lb_loadbalancer.app.id
  # This is unnecessary in real life.
  # This is required here to let the example work with loadbalancer resource example. 
  depends_on = [vkcs_lb_loadbalancer.app]
}

data "vkcs_networking_port" "app_port" {
  id = data.vkcs_lb_loadbalancer.app.vip_port_id
}

output "used_vips" {
  value       = data.vkcs_networking_port.app_port.all_fixed_ips
  description = "IP addresses of the app"
}
