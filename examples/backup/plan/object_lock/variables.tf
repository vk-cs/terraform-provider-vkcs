variable "enable_object_lock" {
  type        = bool
  default     = false
  description = "Set to true to create a protected backup plan."
}

variable "instance_id" {
  type        = string
  default     = ""
  description = "ID of a VM that supports backup protection; required when enable_object_lock is true."
}
