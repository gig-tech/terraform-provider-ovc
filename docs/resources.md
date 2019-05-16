# Resources

* [ovc_machine](#Resource:-ovc_machine)
* [ovc_disk](#Resource:-ovc_disk)
* [ovc_port_forwarding](#Resource:-ovc_port_forwarding)
* [ovc_cloudspace](#Resource:-ovc_cloudspace)

## Resource: ovc_machine

Provides a ovc machine. This allows machines to be created, updated and deleted

### Example Usage

```hcl
variable "cloudspace_id" {}

resource "ovc_machine" "machine" {
  cloudspace_id = "${var.cloudspace_id}"
  image_id = 5
  size_id = 1
  disksize = 10
  name = "MyMachine"
  description = "Machine Provisioned With Terraform"
}
```

The cloudspace can be obtained with the following commands:

```
export CLIENT_SECRET=
export CLIENT_ID=
export cloudspace_name=example_name
JWT=$(curl -d 'grant_type=client_credentials&client_id='"$CLIENT_ID"'&client_secret='"$CLIENT_SECRET"'&response_type=id_token' https://itsyou.online/v1/oauth/access_token)
export cloudspaces_json=$(curl -X POST -H "Authorization: bearer $JWT" https://ch-lug-dc01-001.gig.tech/restmachine//cloudapi/cloudspaces/list)
echo $cloudspaces_json | jq -r 'map(select(any(.name; contains($cn)))|.id)[]' --arg cn "$cloudspace_name"

```

### Argument Reference

The following arguments are supported:

* cloudspace_id - (Required) The cloudspace ID of the cloudspace where the machine needs to be created
* image_id - (Required) The image ID of the image to use for this instance
* size_id - (Required) Size ID for this instance
* disksize - (Required) Size of the boot disk in gigabytes
* iops - (Optional) IOPS limiting of the boot disk
* name - (Required) Name of the machine
* description - (Optional) Description of the machine

## Resource: ovc_disk

Creates extra disks used by ovc machines

### Example Usage

```hcl
variable "cloudspace_id" {}

resource "ovc_machine" "machine" {
  cloudspace_id = "${var.cloudspace_id}"
  image_id = 5
  size_id = 1
  disksize = 10
  name = "MyMachine"
  description = "Machine Provisioned With Terraform"
}

resource "ovc_disk" "disk1" {
  machine_id = "${ovc_machine.machine.id}"
  disk_name = "terraform_disk"
  description = "Disk created by terraform"
  size = 10
  type = "D"
  iops = 1000
}

resource "ovc_disk" "disk2" {
  machine_id = "${ovc_machine.machine.id}"
  disk_name = "terraform_disk"
  description = "Disk created by terraform"
  size = 15
  type = "D"
  iops = 1500
  depends_on = ["ovc_disk.disk1"]
}
```

For creating several disks add resource construct for each disk separately, add dependency between disks to create them sequentially.
Terraform allows creating multiple objects with attribute `count`, but in this case it is not allowed, since OVC supports only adding one disks to a VM at a time.

### Argument Reference

The following arguments are supported:

* machine_id - (Required) Machine ID of the machine where the disk should be attached
* disk_name - (Required) Disk name of the disk
* description - (Required) Disk description
* size - (Required) Size in gigabytes of the disk
* type - (Required) Type of disk, following options are supported: B (Boot), D (Data)
* iops - (Optional) Maximum IOPS disk can perform, defaults to 2000

## Resource: ovc_port_forwarding

Manages port forwarding

### Example Usage

```hcl
variable "cloudspace_id" {}
variable "cloudspace_public_ip" {}

resource "ovc_machine" "machine" {
  cloudspace_id = "${var.cloudspace_id}"
  image_id = 5
  size_id = 1
  disksize = 10
  name = "MyMachine"
  description = "Machine Provisioned With Terraform"
}

resource "ovc_port_forwarding" "port_forward" {
  cloudspace_id = "${var.cloudspace_id}"
  public_ip = "${var.cloudspace_public_ip}"
  public_port = 222
  machine_id = "${ovc_machine.machine.id}"
  local_port = 22
  protocol = "tcp"
}
```

The public ip can be obtained via the API:

```
export CLIENT_SECRET=
export CLIENT_ID=
export cloudspace_name=example_name
JWT=$(curl -d 'grant_type=client_credentials&client_id='"$CLIENT_ID"'&client_secret='"$CLIENT_SECRET"'&response_type=id_token' https://itsyou.online/v1/oauth/access_token)
export cloudspaces_json=$(curl -X POST -H "Authorization: bearer $JWT" https://ch-lug-dc01-001.gig.tech/restmachine//cloudapi/cloudspaces/list)
echo $cloudspaces_json | jq -r 'map(select(any(.name; contains($cn)))|.externalnetworkip)[]' --arg cn "$cloudspace_name"
```

### Argument Reference

* cloudspace_id - (Required) ID of the cloudspace
* public_ip - (Required) public ip of the cloudspace
* public_port - (Required) public port which should be forwarded
* machine_id - (Required) machine ID of where to forward the port to
* local_port - (Required) local port of the machine where to forward to
* protocol - (Required) protocol to use, either "tcp" or "udp"

## Resource: ovc_cloudspace

Creates cloudpsaces

### Example Usage

```hcl
resource "ovc_cloudspace" "cloudspace" {
  account = "${var.account}"
  name = "cloudspace"
  private_network = "192.168.100.0/24"
   resource_limits = {
     max_memory_capacity = 3.0
     max_disk_capacity = 12
     max_cpu_capacity = 3
     max_num_public_ip = 4
     max_network_peer_transfer = 4
   }
}
```

To get the account ID you can query the API:

```
export CLIENT_SECRET=
export CLIENT_ID=
export cloudspace_name=
export account_name=
JWT=$(curl -d 'grant_type=client_credentials&client_id='"$CLIENT_ID"'&client_secret='"$CLIENT_SECRET"'&response_type=id_token' https://itsyou.online/v1/oauth/access_token)
export accounts_json=$(curl -X POST -H "Authorization: bearer $JWT" https://ch-lug-dc01-001.gig.tech/restmachine//cloudapi/accounts/list)
echo $accounts_json | jq -r 'map(select(any(.name; contains($account_name)))|.id)[]' --arg account_name "$account_name"
```

### Argument Reference

* `account` - (Required) Name of the account this cloudspace belongs to
* `name` - (Required) name of space to create
* `private_network` - (Optional) private network CIDR eg. 192.168.103.0/24
* `resource_limits` - (Optional) specify resource limits block
  * `max_memory_capacity` - (Optional) max size of memory in GB
  * `max_disk_capacity` - (Optional) max size of aggregated vdisks in GB
  * `max_cpu_capacity` - (Optional) max number of cpu cores
  * `max_num_public_ip` - (Optional) max number of assigned public IPs
  * `max_network_peer_transfer` - (Optional) max sent/received network transfer peering

## Resource: ovc_ipsec

Manages port forwarding

### Example Usage

```hcl
provider "ovc" {
  server_url = "${var.server_url}"
  client_jwt="${var.client_jwt}"
}
resource "ovc_cloudspace" "cs" {
  account = "${var.account}"
  name = "${var.cloudspace1}"
  private_network = "192.168.103.0/24"
}
resource "ovc_cloudspace" "cst" {
  account = "${var.account}"
  name = "${var.cloudspace2}"
  private_network = "192.168.104.0/24"
}
data "ovc_cloudspace" "cs" {
  account = "${var.account}"
  name = "${var.cloudspace1}"
  depends_on = ["ovc_cloudspace.cs"]
}
data "ovc_cloudspace" "cst" {
  account = "${var.account}"
  name = "${var.cloudspace2}"
  depends_on = ["ovc_cloudspace.cst"]
}
resource "ovc_ipsec" "tunnel1" {
  cloudspace_id = "${ovc_cloudspace.cs.id}"
  remote_public_ip = "${data.ovc_cloudspace.cst.external_network_ip}"
  remote_private_network = "${data.ovc_cloudspace.cst.private_network}"
  depends_on = ["ovc_cloudspace.cs", "ovc_cloudspace.cst"]
}
resource "ovc_ipsec" "tunnel2" {
  cloudspace_id = "${ovc_cloudspace.cst.id}"
  remote_public_ip = "${data.ovc_cloudspace.cs.external_network_ip}"
  remote_private_network = "${data.ovc_cloudspace.cs.private_network}"
  psk = "${ovc_ipsec.tunnel1.psk}"
  depends_on = ["ovc_cloudspace.cs", "ovc_cloudspace.cst", "ovc_ipsec.tunnel1"]

```

### Argument Reference

* cloudspace_id - (Required) ID of the cloudspace
* remote_public_ip - (Required) public ip of the cloudspace to connect to
* remote_private_network - (Required) remote private network to connect to
* psk - (Optional) Pre shared secret for the connection's authentication
