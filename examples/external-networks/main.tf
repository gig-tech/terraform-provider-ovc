provider "ovc" {
  server_url = "${var.server_url}"
}

// use this data source to get list all network available for this account

data "ovc_external_networks" "nets"{
  account = "{}"
}

// use this data source to get external network by name
data "ovc_external_network" "net"{
  name = "${var.external_network}"
}

data "ovc_cloudspace" "cs"{
  account = "${var.account}"
  name = "${var.cs_name}"
}

# Data definition for image
data "ovc_image" "im"{
  most_recent = true
  name_regex = "(?i).*\\.?ubuntu.*16*"
}

resource "ovc_machine" "machine" {
  cloudspace_id = "${data.ovc_cloudspace.cs.id}"
  image_id      = "${data.ovc_image.ubuntu16.image_id}"
  size_id       = "${var.size_id}"
  disksize      = "${var.disksize}"
  name          = "${var.machine}"
  description   = "${var.vm_description}"
  interfaces    = [
    # if empty interface is given, VM will be attached to the default external network
    {}, 
    # if network ID is given, VM will be attached to this network
    {network_id = "${data.ovc_external_network.net.id}"}, 
    # Several IP addresses from the same external network can be added, just add nics for the same network ID
    {network_id = "${data.ovc_external_network.net.id}"}
  ]
}
// print VM id
output "id"{
  value = "${ovc_machine.vm.id}"

// list external networks
output "nets"{
  value = "${data.ovc_external_networks.nets.entities}"
}

// list VM external nics
output "attched_external_networks"{
  value = "${ovc_machine.vm.interfaces}"
}

}
