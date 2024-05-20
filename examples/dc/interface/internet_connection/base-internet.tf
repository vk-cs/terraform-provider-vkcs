# Get external network with Internet access
data "vkcs_networking_network" "internet_sprut" {
  name = "internet"
  sdn  = "sprut"
}
