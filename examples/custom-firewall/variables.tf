variable "server_url"{
  description = "API server URL"
}

variable "client_jwt" {
  description = "itsyouonline jwt token"
}

variable "account" {
  description = "account"
}

variable "disksize" {
  description = "data disk size"
  default = 10
}

variable "external_network" {
  description = "external network name"
  description = ""
}
variable "vm_description" {
  description = "machine description"
  default = "machine deployed with Terraform"
}
variable "size_id" {
  description = "machine dize id"
  default = 3
}

variable "userdata" {
  description = "user info"
  default = "users: [{name: root, shell: /bin/bash, ssh-authorized-keys: []}]"
}
