variable "name" {
  type = string
  default = "default"
}

variable "max-node-unavailable" {
  type = number
  default = 1
}

variable "dns-domain" {
  type = string
  default = "cluster.local"
}
