# OVC Account name, your IYO account must have access to it.
account = "<Account Name>"

# IYO client id and client secret.
client_id = ""
client_secret = ""

# G8 api url
server_url= "<Server URL>"

# cloudspace name
cs_name = "<CS Name>"

# image id of the image that will be used to create virtual machines
# Needs to be looked up through the admin interface or fetched from ovc_image data source
# image_id = 286

# Needs to be looked up through the API!
# Size id (define how many cpus / ram )  you can get that from ovc first
size_id = 3

# disk size in GB you can get that from ovc first
disksize = 10

# The description of the VM
vm_description = "Meneja K8S"

# User data, contain users and SSH keys to be added to the VM
userdata = "users: [{name: root, shell: /bin/bash, ssh-authorized-keys: [key1, key2]}]"