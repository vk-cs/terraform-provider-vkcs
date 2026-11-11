variable "enable_restore" {
  type        = bool
  default     = false
  description = "Set to true to restore a VM from the selected checkpoint."
}

variable "instance_name" {
  type        = string
  default     = "restored-vm-example"
  description = "Name for the restored VM."
}

variable "checkpoint_id" {
  type        = string
  default     = ""
  description = "ID of the backup checkpoint to restore."
}

variable "source_resource_id" {
  type        = string
  default     = ""
  description = "Nova server resource ID within the checkpoint when it contains multiple VMs."
}

variable "network_id" {
  type        = string
  default     = ""
  description = "Network ID for the restored VM."
}

variable "key_pair" {
  type        = string
  default     = ""
  description = "Keypair name for the restored VM."
}

variable "flavor_id" {
  type        = string
  default     = ""
  description = "Optional flavor ID for the restored VM; defaults to the flavor in the checkpoint."
}
