# OVC Account name, your IYO account must have access to it.
account = "<Account Name>"

# G8 api url
server_url= "<Server URL>"

disksize = "<Disk size required for checkpoint image>"
external_network = "<External Newtork Name>"

# User data to be added to the VM 
userdata = "users: [{name: user, shell: /bin/bash, ssh-authorized-keys:[]},{name: root, shell: /bin/bash, ssh-authorized-keys: [<SSH KEYS to add on machine>]}]"
