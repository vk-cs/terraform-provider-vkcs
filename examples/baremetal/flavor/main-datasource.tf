data "vkcs_baremetal_flavor" "main" {
  name              = "BM_CX301_N_BOND"
  cpu_model         = "Intel(R) Xeon(R) Gold 6338 CPU @ 2.00GHz"
  cpu_cores         = 32
  ram_size          = 128
  ssd_size          = 900
  hdd_size          = 16000
  bond_vlan_capable = true
}

output "flavor_output" {
  value = {
    id                = data.vkcs_baremetal_flavor.main.id
    name              = data.vkcs_baremetal_flavor.main.name
    cpu_cores         = data.vkcs_baremetal_flavor.main.cpu_cores
    ram_size          = data.vkcs_baremetal_flavor.main.ram_size
    ssd_size          = data.vkcs_baremetal_flavor.main.ssd_size
    hdd_size          = data.vkcs_baremetal_flavor.main.hdd_size
    bond_vlan_capable = data.vkcs_baremetal_flavor.main.bond_vlan_capable
  }
}
