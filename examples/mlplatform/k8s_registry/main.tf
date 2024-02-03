resource "vkcs_mlplatform_k8s_registry" "k8s_registry" {
  name              = "tf-example"
  admin_name        = "admin"
  admin_password    = "dM8Ao21,0S264iZp"
  flavor_id         = data.vkcs_compute_flavor.basic.id
  availability_zone = "GZ1"
  boot_volume = {
    volume_type = "ceph-ssd"
  }
  networks = [
    {
      network_id = vkcs_networking_network.app.id
      ip_pool = data.vkcs_networking_network.extnet.id
    },
  ]
}
