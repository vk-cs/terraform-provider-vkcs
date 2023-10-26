resource "vkcs_networking_subnet" "subnet-with-dns-tf-example" {
  name       = "subnet-with-dns-tf-example"
  network_id = vkcs_networking_network.app.id
  # here we set DNS'es instead of built in ones.
  # subnet resources will not be available via their private DNS names.
  dns_nameservers = [
    "8.8.8.8",
    "8.8.4.4"
  ]
  cidr = "192.168.200.0/24"
}
