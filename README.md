# OpenvCloud Provider

The OpenvCloud provider is used for interacting with the cloud platform in order to create the desired resources. The provider needs to be configured with credentials.

## Build

To build, make sure the project is in your GOPATH and run the following command:


```
mkdir -p ~/.terraform.d/plugins
mkdir -p $GOPATH/src/github.com/terraform-providers/terraform-provider-ovc
git clone git@git.gig.tech:nubera/terraform-provider-ovc.git $GOPATH/src/github.com/terraform-providers/terraform-provider-ovc
go get -v -u github.com/hashicorp/terraform/terraform
go get -v -u github.com/nuberabe/ovc-sdk-go/ovc
cd $GOPATH/src/github.com/terraform-providers/terraform-provider-ovc
go build -o terraform-provider-ovc .
mv terraform-provider-ovc ~/.terraform.d/plugins
```

Put the binary in your plugins folder of terraform. For more information:

https://www.terraform.io/docs/plugins/basics.html#installing-plugins

## Example Usage

```hcl
variable "client_id" {}
variable "client_secret" {}
variable "ovc_url" {}

# Configure the ovc provider
provider "ovc" {
  server_url = "${var.ovc_url}"
  client_id = "${var.client_id}"
  client_secret = "${var.client_secret}"
}
```

## Authentication

Authentication is through itsyouonline with a client ID and client secret


## Argument Reference

The following arguments are supported in the provider block:

* server_url - (Required) The server url of the ovc api to connect to
* client_id - (Required) The client_id of the itsyouonline user
* client_secret - (Required) The client_secret of the itsyouonline user

The arguments can be provided as environment variables:

```
export OPENVCLOUD_SERVER_URL="server-url"
export ITSYOU_ONLINE_CLIENT_ID="your-client-id"
export ITSYOU_ONLINE_CLIENT_SECRET="your-client-secret"
```
This way the provider information must not be included in your terraform configuration file.

# Resources

* [ovc_machine](#Resource:-ovc_machine)
* [ovc_disk](#Resource:-ovc_disk)
* [ovc_port_forwarding](#Resource:-ovc_port_forwarding)
* [ovc_cloudspace](#Resource:-ovc_cloudspace)

# Data Sources

* [ovc_machine](#Data-source:-ovc_machine)
* [ovc_cloudspace](#Data-source:-ovc_cloudspace)
* [ovc_disk](#Data-source:-ovc_disk)
* [ovc_sizes](#Data-source:-ovc_sizes)
* [ovc_machines](#Data-source:-ovc_machines)
* [ovc_cloudspaces](#Data-source:-ovc_cloudspaces)

# Resource: ovc_machine

Provides a ovc machine. This allows machines to be created, updated and deleted

## Example Usage

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

## Argument Reference

The following arguments are supported:

* cloudspace_id - (Required) The cloudspace ID of the cloudspace where the machine needs to be created
* image_id - (Required) The image ID of the image to use for this instance
* size_id - (Required) Size ID for this instance
* disksize - (Required) Size of the boot disk in gigabytes
* name - (Required) Name of the machine
* description - (Optional) Description of the machine


# Resource: ovc_disk

Creates extra disks used by ovc machines

## Example Usage

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

resource "ovc_disk" "disk" {
  machine_id = "${ovc_machine.machine.id}"
  disk_name = "terraform_disk"
  description = "Disk created by terraform"
  size = 10
  type = "D"
  ssd_size = 10
  iops = 2000
}
```

## Argument Reference

The following arguments are supported:

* machine_id - (Required) Machine ID of the machine where the disk should be attached
* disk_name - (Required) Disk name of the disk
* description - (Required) Disk description
* size - (Required) Size in gigabytes of the disk
* type - (Required) Type of disk, following options are supported: B (Boot), D (Data), T (Temp)
* ssdSize - (Optional) Size in gigabytes of the ssd disk
* iops - (Optional) Maximum IOPS disk can perform, defaults to 2000

# Resource: ovc_port_forwarding

Manages port forwarding

## Example Usage

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

## Argument Reference

* cloudspace_id - (Required) ID of the cloudspace
* public_ip - (Required) public ip of the cloudspace
* public_port - (Required) public port which should be forwarded
* machine_id - (Required) machine ID of where to forward the port to
* local_port - (Required) local port of the machine where to forward to
* protocol - (Required) protocol to use, either "tcp" or "udp"

# Resource: ovc_cloudspace

Creates cloudpsaces

## Example Usage

```hcl
resource "ovc_cloudspace" "cloudspace" {
  account = "${var.account}"
  name = "cloudspace"
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

## Argument Reference

* account - (Required) Name of the account this cloudspace belongs to
* name - (Required) name of space to create
* resource_limits - (Optional) specify resource limits block 
  *  max_memory_capacity - (Optional) max size of memory in GB
  *  max_disk_capacity - (Optional) max size of aggregated vdisks in GB
  *  max_cpu_capacity - (Optional) max number of cpu cores
  *  max_num_public_ip - (Optional) max number of assigned public IPs
  *  max_network_peer_transfer - (Optional) max sent/received network transfer peering


# Data Source: ovc_machine

Use this data source to get the ID of a machine in a cloudspace by name

## Example Usage

```hcl
data "ovc_machine" "machine" {
   cloudspace_id = "${var.cloudspace_id}"
   name = "${var.name}"
}
```

## Argument Reference

* name - (Required) name of machine to look up
* cloudspace_id - (Required) ID of the cloudspace where the machine is located

# Data Source: ovc_cloudspace

Use this data source to get the ID of a cloudspace in a location by name

## Example Usage

```hcl
data "ovc_cloudspace" "cloudspace" {
  account = "${var.accountname}"
  name = "${var.name}"
}
```

## Argument Reference

* name - (Required) name of cloudspace to look up
* account - (Required) name of the account where the cloudspace is located

# Data Source: ovc_sizes

Use this data source to get ID of sizes given vcpus and memory

## Example Usage

```hcl
data "ovc_sizes" "size" {
  cloudspace_id = 225
  vcpus = 2
  memory = 4096
}
```

## Argument Reference

* cloudspace_id - (Required) cloudspace ID where the size is located
* memory - (Required) memory of the size
* vcpus - (Required) vcpus of the size 

# Data Source: ovc_disk

Use this data source to get the ID of a disk in a cloudspace by name

## Example Usage

```hcl
data "ovc_disk" "disk" {
   account_id = 4
   name = "${var.name}"
}
```

## Argument Reference

* account_id - (Required) ID of the account where the disk is located
* name - (Required) name of the disk to look up

# Data Source: ovc_machines

Use this data source to retrieve information about all machines in a given cloudspace

## Example Usage

```hcl
data "ovc_machines" "machines" {
   cloudspace_id = "${var.cloudspace_id}"
}
```

## Argument Reference

* cloudspace_id - (Required) ID of the cloudspace where the machines are located

# Data Source: ovc_cloudspaces

Use this data source to retrieve information about all cloudspaces in a given location

## Example Usage

```hcl
data "ovc_cloudspaces" "cloudspaces" {
}
```

## Argument Reference

* No arguments needed for this data source
