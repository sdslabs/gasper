# Iwa Configuration

Iwa service provides [SSH](https://www.ssh.com/ssh/protocol/) access directly to an application's docker container to the end user

The SSH command will be automatically returned to the user on application creation provided the node where the application is deployed has the Iwa service deployed

The following section deals with the configuration of Iwa

```toml
#########################
#   Iwa Configuration   #
#########################

[services.iwa]
deploy = false   # Deploy Iwa?
port = 2222

# Location of Private Key for creating the SSH Signer.
host_signers = ["/home/user/.ssh/id_rsa"]
using_passphrase = false   # Private Key is passphrase protected?
passphrase = ""   # Passphrase (if any) for decrypting the Private Key

# IP address to establish a SSH connection to.
# Equal to the current node's IP address if left blank.
# This field is only for information of the client who will create applications 
# and this field's value will not affect Iwa's functioning in any manner.
# To be used when the current node is only accessible by a jump host or 
# behind some network forwarding rule or proxy setup.
entrypoint_ip = ""
```

The **host_signers** field stores the location of your private key

!!!info
    If your private key is passphrase protected then set the **using_passphrase** field to `true` and insert your passphrase as the value of the **passphrase** field

!!!info
    The password required for SSH access is provided by the user during application creation 
