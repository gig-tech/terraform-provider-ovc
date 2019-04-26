# Data Sources

* [ovc_machine](#Data-source:-ovc_machine)
* [ovc_cloudspace](#Data-source:-ovc_cloudspace)
* [ovc_disk](#Data-source:-ovc_disk)
* [ovc_sizes](#Data-source:-ovc_sizes)
* [ovc_machines](#Data-source:-ovc_machines)
* [ovc_cloudspaces](#Data-source:-ovc_cloudspaces)
* [ovc_image](#Data-source:-ovc_image)
* [ovc_images](#Data-source:-ovc_images)

## Data Source: ovc_machine

Use this data source to get the ID of a machine in a cloudspace by name

### Example Usage

```hcl
data "ovc_machine" "machine" {
   cloudspace_id = "${var.cloudspace_id}"
   name = "${var.name}"
}
```

### Argument Reference

* name - (Required) name of machine to look up
* cloudspace_id - (Required) ID of the cloudspace where the machine is located

## Data Source: ovc_cloudspace

Use this data source to get the ID of a cloudspace in a location by name

### Example Usage

```hcl
data "ovc_cloudspace" "cloudspace" {
  account = "${var.accountname}"
  name = "${var.name}"
}
```

### Argument Reference

* name - (Required) name of cloudspace to look up
* account - (Required) name of the account where the cloudspace is located

## Data Source: ovc_sizes

Use this data source to get ID of sizes given vcpus and memory

### Example Usage

```hcl
data "ovc_sizes" "size" {
  cloudspace_id = 225
  vcpus = 2
  memory = 4096
}
```

### Argument Reference

* cloudspace_id - (Required) cloudspace ID where the size is located
* memory - (Required) memory of the size
* vcpus - (Required) vcpus of the size 

# Data Source: ovc_disk

Use this data source to get the ID of a disk in a cloudspace by name

### Example Usage

```hcl
data "ovc_disk" "disk" {
   account_id = 4
   name = "${var.name}"
}
```

### Argument Reference

* account_id - (Required) ID of the account where the disk is located
* name - (Required) name of the disk to look up

## Data Source: ovc_machines

Use this data source to retrieve information about all machines in a given cloudspace

## Example Usage

```hcl
data "ovc_machines" "machines" {
   cloudspace_id = "${var.cloudspace_id}"
}
```

### Argument Reference

* cloudspace_id - (Required) ID of the cloudspace where the machines are located

## Data Source: ovc_cloudspaces

Use this data source to retrieve information about all cloudspaces in a given location

### Example Usage

```hcl
data "ovc_cloudspaces" "cloudspaces" {
}
```

### Argument Reference

* No arguments needed for this data source


## Data Source: ovc_image

Use this data source to retrieve image by name. If more that single image matches the query, Terraform will fail.
To return the list of images use the data source [`ovc_images`](#Data-Source:-ovc_images).

### Example Usage

```hcl
data "ovc_image" "im"{
  account = "<Account Name>"
  name_regex = "<ImageName>"
  most_recent = true
}
```

### Argument Reference

* `account` - (Optional) name of the account to retrieve images from. If set to 0, only system images will be looked up.
* `name_regex` - (Optional) full name or name pattern for regex search. If set to "" all available images will be looked up
* `most_recent` - (Optional) If set to `true` will search for the latest crated image within the scope (image with the largest ID)

## Data Source: ovc_images

Use this data source to retrieve list of images by name

### Example Usage

```hcl
data "ovc_image" "im"{
  account = "<Account Name>"
  name_regex = "<ImageName>"
}
```

### Argument Reference

* `account` - (Optional) name of the account to retrieve images from. If set to 0, only system images will be looked up.
* `name_regex` - (Optional) full name or name pattern for regex search. If set to "" all available images will be looked up
