variable "enable_checkpoints" {
  type        = bool
  default     = false
  description = "Set to true to query checkpoints from an existing backup plan."
}

variable "backup_plan_id" {
  type        = string
  default     = ""
  description = "ID of an existing backup plan; used when enable_checkpoints is true."
}
