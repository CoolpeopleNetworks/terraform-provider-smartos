SmartOS Cluster Terraform Provider
=========================

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

NOTE: This is a preliminary checkin of a SmartOS provider that supports talking to multiple nodes in a cluster.  This is a fork (but not in the git sense) of my SmartOS provider.   The below documentation is not yet updated to reflect this provider's new capabilities.


Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.10.x
-	[Go](https://golang.org/doc/install) 1.9 (to build the provider plugin)

Using the provider
------------------

This provider can be used to provision machines with a SmartOS host via SSH.  SSH public keys are expected to already be installed on the SmartOS host in order for this provider to work.   

NOTE: Currently, this provider only supports a subset of properties for SmartOS virtual machines.

### Setup ###

```hcl
provider "smartos" {
    hosts = {
            "node01" = "10.99.50.60"
            }
    user = "root"
    sshkey = "~/.ssh/id_ed25519"
}
```

The following arguments are supported.

- `hosts` - (Required) These are the addresses of the global zone on the SmartOS hosts.
- `user` - (Required) This is the authenticated SSH user which will run provisioning commands.   Normally this is 'root'.
- `sshkey` - (Required) This is the private SSH key of the provisioning user. Only passwordless keys work at the moment.

### Resources and Data Providers ###

Currently, the following data and resources are provided:

- smartos_image (Data source - the images will be imported on first use by a smartos_machine stanza.)
- smartos_machine (Resource)

NOTE: The property names supported by this provider match (as much as possible) those defined by Joyent for use with their 'vmadm' utility.   See the man page (specifically the PROPERTIES section) for that utility for more info:

https://smartos.org/man/1m/vmadm

Many of the properties defined in the man page are not yet supported by the provider.

### Example ###

The following example shows you how to configure a bhyve VM running provided Debian 11 image.

```hcl
provider "smartos" {
    hosts = {
          "node01" = "10.99.50.60"
         }
    user = "root"
    sshkey = "~/.ssh/id_ed25519"
}

data "smartos_image" "bhyve_debian11" {
    node_name = "node01"
    name = "debian-11"
    version = "20220228"
}

resource "smartos_machine" "linux-byve" {
    node_name = "node01"
    alias = "provider-test-linux-bhyve"
    brand = "bhyve"
    vcpus = 2

    customer_metadata = {
        "root_authorized_keys" = "ssh-ed25519 AAAA......."
    }

    maintain_resolvers = true
    ram = 512
    nics {
            nic_tag = "admin"
            ips = ["192.168.0.10/24"]
            gateways = ["192.168.0.1"]
            vlan_id = "10"
            interface = "net0"
            model = "virtio"
        }
    quota = 25

    resolvers = ["1.0.0.1", "1.1.1.1"]

    disks {
            boot = true
            image_uuid = "${data.smartos_image.bhyve_debian11.id}"
            compression = "lz4"
            model = "virtio"
        }

    provisioner "remote-exec" {
        inline = [
            "apt-get update",
            "apt-get -y install htop",
        ]
    }
}



```

Links:
https://learn.hashicorp.com/tutorials/terraform/provider-release-publish?utm_source=WEBSITE&utm_medium=WEB_IO&utm_offer=ARTICLE_PAGE&utm_content=DOCS
