provider "ovc" {
  server_url = "${var.server_url}"
  client_jwt = "${var.client_jwt}"
}
# Resource definition for private cloudspace with no public interface
resource "ovc_cloudspace" "private" {
  account = "${var.account}"
  name = "PRIVATE-CS-TERRAFORM"
  mode = "private"
}
# Resource definition for nested cloudspace
resource "ovc_cloudspace" "nested" {
  account = "${var.account}"
  name = "NESTED-CS-TERRAFORM"
  mode = "nested"
  external_network_id = "${ovc_cloudspace.private.id}" # parent network ID
}
data "ovc_image" "im"{
  most_recent = true
  name_regex = "(?i).*\\.?ubuntu.*16*"
}
data "ovc_image" "checkpoint"{
  most_recent = true
  name_regex = "Checkpoint" # set name of chackpoint image uploaded to OVC
}
data "ovc_external_network" "net"{
  name = "${var.external_network}"
}
# this is castomizable firewall
resource "ovc_machine" "checkpoint" {
  cloudspace_id = "${ovc_cloudspace.private.id}"
  image_id      = "${data.ovc_image.checkpoint.image_id}"
  size_id       = "${var.size_id}"
  disksize      = "${var.disksize}"
  name          = "CHECKPOINT-VM"
  description   = "${var.vm_description}"
  userdata      = "${var.userdata}"
  act_as_default_gateway = true
  interfaces = [
    {
      "network_id" = "${data.ovc_external_network.net.id}"
    }
  ] 
}
# this machine will have access to the internes via the checkpoint
resource "ovc_machine" "server" {
  cloudspace_id = "${ovc_cloudspace.nested.id}"
  image_id      = "${data.ovc_image.im.image_id}"
  size_id       = "${var.size_id}"
  disksize      = "${var.disksize}"
  name          = "SERVER-VM"
  description   = "${var.vm_description}"
  userdata      = "${var.userdata}"
}
