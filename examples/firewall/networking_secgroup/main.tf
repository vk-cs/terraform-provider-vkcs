resource "vkcs_networking_secgroup" "etcd" {
  name        = "etcd-tf-example"
  description = "etcd service"
}
