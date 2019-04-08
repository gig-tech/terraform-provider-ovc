variable "client_id" {
  description = "Client ID created on itsyou.online"
}

variable "client_secret" {
  description = "Client secret created on itsyou.online"
}

variable "server_url" {
  description = "API server URL"
}

variable "account" {
  description = "account"
}

variable "cs_name" {
  description = "cloudspace name"
}

variable "vm_description" {
  description = "Description of the VM"
}

variable "image_id" {
  description = "Image_id"
  default     = "1"
}

variable "size_id" {
  description = "size_id"
  default     = "3"
}

variable "disksize" {
  description = "disksize"
  default     = "20"
}

variable "userdata" {
  description = "user data"
  default = ""
}
