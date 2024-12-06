data "vkcs_networking_secgroup" "etcd" {
  name = "etcd-tf-example"
  # This is unnecessary in real life.
  # This is required here to let the example work with secgroup resource example. 
  depends_on = [vkcs_networking_secgroup.etcd]
}
