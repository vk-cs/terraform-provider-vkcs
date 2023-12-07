resource "vkcs_mlplatform_mlflow_deploy" "deploy" {
  name               = "tf-example"
  flavor_id          = data.vkcs_compute_flavor.basic.id
  mlflow_instance_id = vkcs_mlplatform_mlflow.mlflow.id
  availability_zone  = "GZ1"
  boot_volume = {
    size        = 50
    volume_type = "ceph-ssd"
  }

  data_volumes = [
    {
      size        = 60
      volume_type = "ceph-ssd"
    },
  ]
  networks = [
    {
      network_id = vkcs_networking_network.app.id
    },
  ]
}