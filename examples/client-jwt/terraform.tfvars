# OVC Account name, your IYO account must have access to it.
account = "<Account Name>"

# IYO client id and client secret.
client_jwt = "<IYO JWT>"

# G8 api url
server_url= "<Server URL>"

# cloudspace name
cs_name = "<CS Name>"

# Needs to be looked up through the API!
# Size id (define how many cpus / ram )  you can get that from ovc first
size_id = 3

# disk size in GB you can get that from ovc first
disksize = 10

# The description of the VM
vm_description = "Terraform test VM"

# User data to be added to the VM 
userdata = "users: [{name: Carmichael, shell: /bin/bash, ssh-authorized-keys: [<public key to be added to the VM's authorized_keys>]}]"
