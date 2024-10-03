resource "vkcs_cdn_origin_group" "origin_group" {
  name = "tfexample-origin-group"
  origins = [
    {
      source = "origin1.vk.com"
    },
    {
      source = "origin2.vk.com",
      backup = true
    }
  ]
  use_next = true
}
