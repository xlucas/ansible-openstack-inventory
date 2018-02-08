# ansible-openstack-inventory

A dynamic inventory for openstack clouds with an uncluttered output.

## Installation

```bash
go get github.com/xlucas/ansible-openstack-inventory
```

## Usage

Create a `.clouds.hcl` file in your home directory. This file follows the
[HCL](https://github.com/hashicorp/hcl) format to describe your cloud providers
and their regions.

Below is a sample content for [OVH Cloud](https://www.ovhcloud.com):

```hcl
provider "ovh" {
  options {
    meta {
      user   = "image_original_user" // image metadata used to guess ssh user
      groups = "groups"              // server metadata used to define groups
    }

    fallback_user = "debian" // ssh user to fall back to if metadata is missing
  }

  identity {
    version   = 2
    endpoint  = "https://auth.cloud.ovh.net"
    username  = "xxxxxxxxxxxx"
    password  = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
    tenant_id = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  }

  regions "eu-west" {
    region "Gravelines" {
      name = "GRA1"
    }

    region "Strasbourg" {
      name = "SBG1"
    }

    region "London" {
      name = "UK1"
    }
  }

  regions "eu-central" {
    region "Frankfurt" {
      name = "DE1"
    }
  }

  regions "eu-east" {
    region "Warsaw" {
      name = "WAW1"
    }
  }
}
```

Configure ansible to use it as the default inventory in `ansible.cfg`:
```ini
[defaults]
inventory = /path/to/ansible-openstack-inventory
```
