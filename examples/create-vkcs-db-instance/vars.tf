variable "public-key-file" {
  type    = string
  default = "~/.ssh/id_rsa.pub"
}

variable "db-instance-flavor" {
  type    = string
  default = "Basic-1-2-20"
}