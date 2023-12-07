resource "vkcs_mlplatform_jupyterhub" "jupyterhub" {
  name              = "tf-example"
  admin_name        = "admin"
  admin_password    = "Password!"
  flavor_id         = data.vkcs_compute_flavor.basic.id
  availability_zone = "GZ1"
  boot_volume = {
    volume_type = "ceph-ssd"
  }
  data_volumes = [
    {
      size        = 60
      volume_type = "ceph-ssd"
    },
    {
      size        = 70
      volume_type = "ceph-ssd"
    }
  ]
  networks = [
    {
      network_id = vkcs_networking_network.app.id
    },
  ]
}