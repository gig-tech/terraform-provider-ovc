provider "ovc" {
  server_url = "${var.server_url}"
}

# Definition of the our cloudspace
#resource "ovc_cloudspace" "cs" {
#  account = "${var.account}"
#  name = "${var.cs_name}"
#}

# Data definition for every cloudspace
# To be able to get the ip address
data "ovc_cloudspace" "cs" {
  account = "${var.account}"
  name = "${var.cs_name}"
}

# Data definition for image
data "ovc_image" "im"{
  most_recent = true
  name_regex = "(?i).*\\.?ubuntu.*16*"
}

# Definition of the vm to be created with the settings defined in terraform.tfvars
resource "ovc_machine" "mymachine" {
 cloudspace_id = "${data.ovc_cloudspace.cs.id}"
  image_id      = "${data.ovc_image.im.image_id}"
  size_id       = "${var.size_id}"
  disksize      = "${var.disksize}"
  name          = "mymachine"
  description   = "${var.vm_description}"
  userdata      = "${var.userdata}"
}

resource "ovc_port_forwarding" "ssh" {

  cloudspace_id = "${data.ovc_cloudspace.cs.id}"
  public_ip     = "${data.ovc_cloudspace.cs.external_network_ip}"
  public_port   = 2222
  machine_id    = "${ovc_machine.mymachine.id}"
  local_port    = 22
  protocol      = "tcp"
}
