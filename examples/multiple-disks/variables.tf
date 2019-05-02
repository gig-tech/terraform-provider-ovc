variable "client_jwt" {
  description = "JWT created with itsyou.online API app ID and secret"
}

variable "server_url" {
  description = "API server URL"
}

variable "account" {
  description = "account name"
}

variable "cloudspace" {
  description = "cloudspace name"
}

variable "machine" {
  description = "machine name"
}

variable "size_id" {
  description = "size_id"
  default     = "3"
}

variable "disksize" {
  description = "disksize"
  default     = "20"
}

variable "vm_description" {
  description = "Description of the VM"
}
