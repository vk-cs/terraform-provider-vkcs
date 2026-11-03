variable "backup_plan_id" {
  type        = string
  description = "ID of a backup plan with at least one available VM checkpoint."
}

variable "restore_net_id" {
  type        = string
  description = "Network or subnet ID for the restored VM."
}

variable "restore_keypair" {
  type        = string
  description = "Keypair name for the restored VM."
}
