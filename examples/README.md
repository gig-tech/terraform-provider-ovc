# Examples

This directory contains a few Terraform examples

## [client-id-secret](./client-id-secret):

This is a basic example that uses an app client ID and secret from [itsyou.online](itsyou.online) for authentication.  
It uses the cloudspace defined by `cs_name` in [terraform.tfvars](./client-id-secret/terraform.tfvars) to create a VM `mymachine`
The image ID is directly set.  
There is a port forward defined for the machine, forwarding port `2222`  on the public IP of the cloudspace to `22` on the VM.


## [client-jwt](./client-jwt):

Is a very similar example as [client-id-secret](#[client-id-secret](./client-id-secret)).  
But now uses a JWT from [itsyou.online](itsyou.online) to authenticate and uses the `ovc_image` data source to fetch the latest `ubuntu 16` image ID from OVC using regex.

There is also a port forward defined for the machine, forwarding port `2222`  on the public IP of the cloudspace to `22` on the VM.

In [terraform.tfvars](./client-jwt/terraform.tfvars), the `userdata` is given where a user `Carmichael` is defined to be created with a public key to be added to `Carmichael`'s authorized_keys file.

## [multiple-disks](./multiple-disks)

This is an example on how to set up a VM with a boot disk and 2 additional data disks attached to the VM.
  
Disk resources do not support the `count` parameter, therefor each disk should be defined on it's own. 

Authentication credentials are not defined in the example, instead this example relies on setting the correct environmental variables described in this projects main README under [Authentication](../README.md#Authentication).
