resource "vkcs_dc_router" "dc_router" {
  availability_zone = "MS1"
  flavor            = "standard"
  name              = "tf-example"
  description       = "tf-example-description"
}
