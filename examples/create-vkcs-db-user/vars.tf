variable "db-instance-flavor" {
  type    = string
  default = "Basic-1-2-20"
}

variable "db-user-password" {
  type      = string
  sensitive = true
}
