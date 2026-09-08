data "vkcs_baremetal_flavor" "selected" {
  name = "BM_CX301_N_BOND"
}

# Get the current project-specific display name.
output "display_name_before" {
  value = data.vkcs_baremetal_flavor.selected.display_name
}

# Write a project-specific display name.
resource "vkcs_baremetal_flavor_display_name" "main" {
  id           = data.vkcs_baremetal_flavor.selected.id
  display_name = "Terraform managed bare metal flavor"
}

# Rewrite the display name by changing the value above and running:
# terraform apply

output "display_name" {
  value = vkcs_baremetal_flavor_display_name.main.display_name
}

# Reset the display name by running:
# terraform destroy
# Destroying this resource clears the project-specific display name.
