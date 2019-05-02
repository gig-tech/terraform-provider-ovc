provider "ovc" {
  server_url = "${var.server_url}"
  client_jwt = "${var.client_jwt}"
}

data "ovc_cloudspace" "cs" {
   account = "${var.account}"
   name = "${var.cloudspace}"
}

data "ovc_image" "ubuntu16"{
  most_recent = true
  name_regex = "(?i).*\\.?ubuntu.*16*"
}

# machine definition
resource "ovc_machine" "machine" {
  cloudspace_id = "${data.ovc_cloudspace.cs.id}"
  image_id      = "${data.ovc_image.ubuntu16.image_id}"
  size_id       = "${var.size_id}"
  disksize      = "${var.disksize}"
  name          = "${var.machine}"
  description   = "${var.vm_description}"
}

# Definition of the the disks
resource "ovc_disk" "disk1" {
  machine_id = "${ovc_machine.machine.id}"
  disk_name = "terraform_disk_1"
  description = "Disk created by terraform"
  size = 10
  type = "D"
  iops = 2000
}

resource "ovc_disk" "disk2" {
  machine_id = "${ovc_machine.machine.id}"
  disk_name = "terraform_disk_2"
  description = "Disk created by terraform"
  size = 20
  type = "D"
  iops = 3000
  depends_on = ["ovc_disk.disk1"]
}
