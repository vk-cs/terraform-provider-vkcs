variable "enable_restore" {
  type        = bool
  default     = false
  description = "Set to true to adopt a VM created by a backup restore operation."
}

variable "instance_name" {
  type        = string
  default     = "restored-vm-example"
  description = "Name of the restored VM; it should match the restore operation's restore_name."
}

variable "restore_id" {
  type        = string
  default     = ""
  description = "ID of an existing backup restore operation that restored exactly one VM."
}
