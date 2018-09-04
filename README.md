# ansible-openstack-inventory

An opinionated dynamic inventory for openstack clouds that integrates environment filtering.

## Installation

```bash
go get -u github.com/xlucas/ansible-openstack-inventory
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
      environment     = "ansible_environment" // server metadata used to define environments
      groups          = "ansible_groups"      // server metadata used to define groups
      hostvars_prefix = "ansible_hostvar_"    // server metadata to be added to hostvars
      user            = "image_original_user" // image metadata used to guess ssh user
    }

    fallback_user = "debian" // ssh user to fall back to if image metadata is missing
  }

  identity {
    version   = 2
    endpoint  = "https://auth.cloud.ovh.net" // the identity version is added by the inventory
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

## Usage with environments
- In the configuration file, make sure you have specified the `env` key in the
  `meta` section. This value represents the key that the dynamic inventory will
look for when building the environment host list.
- Make sure your instances are created with the previous metadata key and the
  environment name as a value. It should be the same name than the folder
containing environment-specific configuration in your ansible setup
(`environments/<env>`).
- Make sure the inventory for each environment is invoked with `--env=<env>`.

### Example

Ansible layout:

```text
.
├── ansible.cfg
├── hosts
├── environments
│   ├── dev
│   │   ├── group_vars
│   │   │   └── ...
│   │   └── hosts -> ../../hosts
│   ├── preprod
│   │   ├── group_vars
│   │   │   └── ...
│   │   └── hosts -> ../../hosts
│   └── prod
│       ├── group_vars
│       │   └── ...
│       └── hosts -> ../../hosts
├── roles
│   └── ...
└── playbook.yml
```

The `hosts` file at the root of the directory is a simple script that is used
to invoke the dynamic inventory with the desired environment. Other inventories
are simple symlinks to this script. Its content is:

```bash
#!/bin/sh
dir=${0%/*}
env=${dir##*/}

ansible-openstack-inventory --env=$env ${@}
```

In `ansible.cfg`, specify the default inventory that should be used when not
specified on the command line:

```ini
[defaults]
inventory = ./environments/dev
```

Check that only instances from the development environment are targeted when no
inventory has been specified:
```
ansible -m ping all
4b7ae578-bf37-4179-b324-9fad1b1d78ad | SUCCESS => {
    "changed": false,
    "ping": "pong"
}
```

Run the same command for a specific environment, for instance preproduction:
```
ansible -i environments/preprod -m ping all
b7bd2783-0082-4a4e-997d-00472c6dc940 | SUCCESS => {
    "changed": false,
    "ping": "pong"
}
```

Notice how these hosts are different since they pertain to different
environments.

## Usage without environments

Simply configure ansible to use it as the default inventory in `ansible.cfg`
and omit the `env` parameter in the metadata section.
```ini
[defaults]
inventory = /path/to/ansible-openstack-inventory
```

# License

This project is distributed under the Apache License, Version 2.0.
