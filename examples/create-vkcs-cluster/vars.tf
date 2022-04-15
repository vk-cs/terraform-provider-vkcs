variable "public-key-file" {
    type = string
    default = "~/.ssh/id_rsa.pub"
}

variable "k8s-flavor" {
    type = string
    default = "Standard-2-4-50"
}
