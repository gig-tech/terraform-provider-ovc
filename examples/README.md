# Examples

This directory contains a few Terraform examples. Each example requires several settings to make it work.

Below are the settings that are required in each example.

* `account`, the OVC Account name, your IYO account must have access to it.
* `server_url` G8 url
* `cs_name` the name of the cloudspace

Next to the previous settings, each example can have specific settings. They are mentioned below.


## [client-id-secret](./client-id-secret):

This is a basic example that uses an app client ID and secret from [itsyou.online](itsyou.online) for authentication.
It uses the cloudspace defined by `cs_name` in [terraform.tfvars](./client-id-secret/terraform.tfvars) to create a VM `mymachine`
The image ID is directly set.
There is a port forward defined for the machine, forwarding port `2222`  on the public IP of the cloudspace to `22` on the VM.

Following parameters need to be configured in `terraform.tfvars`:
* `client_id` and `client_secret`

## [client-jwt](./client-jwt):

Is a very similar example as [client-id-secret](#[client-id-secret](./client-id-secret)).
But now uses a JWT from [itsyou.online](itsyou.online) to authenticate and uses the `ovc_image` data source to fetch the latest `ubuntu 16` image ID from OVC using regex.

There is also a port forward defined for the machine, forwarding port `2222`  on the public IP of the cloudspace to `22` on the VM.

In [terraform.tfvars](./client-jwt/terraform.tfvars), the `userdata` is given where a user `Carmichael` is defined to be created with a public key to be added to `Carmichael`'s authorized_keys file.

Following parameters need to be configured in `terraform.tfvars`:
* `client_jwt`

## [multiple-disks](./multiple-disks)

This is an example on how to set up a VM with a boot disk and 2 additional data disks attached to the VM.

Following parameters need to be configured:
* `client_jwt`, as environment variable `TF_VAR_client_jwt`. See [Authentication](../README.md#Authentication).

## [external-networks](./external-networks)

This is an example on how to attach/detach a VM to an external network.

## [custom-firewall](./custom-firewall)

This is an example on how to deploy a cloudspace with custom firewall.