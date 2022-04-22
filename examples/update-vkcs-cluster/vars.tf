variable "public-key-file" {
    type = string
    default = "~/.ssh/id_rsa.pub"
}

variable "k8s-flavor" {
    type = string
    default = "Standard-2-4-50"
}

variable "k8s-network-id" {
    type = string
    default = "95663bae-6763-4a53-9424-831975285cc1"
}

variable "k8s-router-id" {
    type = string
    default = "95663bae-6763-4a53-9424-831975285cc1"
}

variable "k8s-subnet-id" {
    type = string
    default = "95663bae-6763-4a53-9424-831975285cc1"
}

variable "new-master-flavor" {
    type = string
    default = "d659fa16-c7fb-42cf-8a5e-9bcbe80a7538"
}