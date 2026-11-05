variable "enable_restore" {
  type        = bool
  default     = false
  description = "Set to true to restore an available VM checkpoint from the selected plan."
}

variable "backup_plan_id" {
  type        = string
  default     = ""
  description = "ID of a backup plan with at least one available VM checkpoint; used when enable_restore is true."
}

variable "restore_net_id" {
  type        = string
  default     = ""
  description = "Network or subnet ID for the restored VM."
}

variable "restore_keypair" {
  type        = string
  default     = ""
  description = "Keypair name for the restored VM."
}
