data "vkcs_networking_router" "router" {
  name = "router-tf-example"
  tags = ["tf-example"]
  # This is unnecessary in real life.
  # This is required here to let the example work with router resource example. 
  depends_on = [vkcs_networking_router.router]
}
