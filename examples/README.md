# (C) Copyright 2020-2022 Hewlett Packard Enterprise Development LP
# Examples of using the Metal Terraform provider.

## Pre-requisites

Terraform version 0.12.13 and above is required and maybe obtained here `https://www.terraform.io/downloads.html`. 

## Obtaining the provider.

The provider and the portal use a versioned communcations session that must be synchronized. The provider maybe dowloaded from
the portal under the `Downloads` section.

```
  Instructions here on how to down load the provider
```

The provider must be installed where Terraform can locate it. The easiest location is to install the provider 
to `~/.terraform.d/plugins/terraform-provider-hpegl_vx.x.x`.

The provider maybe used with the Metal portal simulator for test purposes. If the test portal is restarted all
simulated resources are lost. In this case, it is important to ensure that all dangling terraform state information
is also delete `rm -f terraform.tfstate*`.

## Installing and setting up the provider

The provider will use the authentication settings of the last successful `steelctl login` command. Infrastructure 
will be reported and manipulated with those credentials, more sepcfiically those found in `~/.qjwt`.

To start using the provider, login to the portal using steelctl, create a new working directory, and write some terraform.

```
# mkdir qform
# cd qform
# terraform init

Initializing the backend...

Initializing provider plugins...

The following providers do not have any version constraints in configuration,
so the latest version was installed.

To prevent automatic upgrades to new major versions that may contain breaking
changes, it is recommended to add version = "..." constraints to the
corresponding provider blocks in configuration, with the constraint strings
suggested below.

* provider.hpegl: version = "~> 0.0"

Terraform has been successfully initialized!

You may now begin working with Terraform. Try running "terraform plan" to see
any changes that are required for your infrastructure. All Terraform commands
should now work.

If you ever set or change modules or backend configuration for Terraform,
rerun this command to reinitialize your working directory. If you forget, other
commands will detect it and remind you to do so if necessary.
# 
```

Terraform scripts may now be created in this directory and then planned or applied using `terraform plan` or `terraform apply` 
respectively. 

Terraform will create terraform.tfstate files in the local directory by default. The actual location of these file maybe changed
using the -state-path option on the command line. These serve as caches of the terraformed infrastructure as described by
the scripts. If changes are made to this infrastructure outside of terraform itself (deleting resources on the portal directly, for instance)
this cache becomes stale. With a stale cache terraform will likely be unable to make the necessary changes without error. Running
`terraform refresh` may go someway to helping resynchronisation but it is not guaranteed.

The corollary is also true; removal of all terraform.tfstate files will leave dangling resources in the portal.


# Using the provider

The terraform provider is used to interact with resources supported by Metal portal services. The provider needs to be
provisioned with credentials for the portal on which resources reside. 

## Example usage
```
provider "hpegl" {
  # A version constraint may be added if later breaking changes to the API force a roll of the provider major number.
  # more details here https://www.terraform.io/docs/configuration/terraform.html
  # version = "~>0.0"  
  metal {
    # It is expected that a user has already authenticated with steelctl login and this terraform 
    # operation will run using the credentials cached in the ~/.qjwt file.

    # portal_url maybe set here and, if set, will be used to verify that the portal used for authentication
    # is the same as that specified here.
    portal_url = "http://172.25.0.2:3002"  # Default simulator address

    # user maybe specified here and, if set, will be used to verify that the authentiction tokens
    # in ~/.qjwt match this user.
    user = "h1@quattronetworks.com"  # a default simulated hoster
  }
}

resource "hpegl_metal_host" "terra_host" {
  name          = "tfhost"
  image         = "CoreOS@2135.6.0"                
  machine_size  = "Any"
  ssh           = ["User1 - Linux"]  
  networks      = ["Private", "Public"]  
  location      = "USA:Texas:AUSL2"
  description   = "Terraformed host"
}

output "tfhost" {
    value = hpegl_metal_host.terra_host.connections
}
```

## Other examples

1. [Available resrources](./available-resources/README.md): Obtain information about unprovisioned resources available to terraform.
1. [Available images](./available-images/README.md): Obtain filtered, specifc image information.
1. [SSH key creation](./ssh-key-creation/README.md): Create new SSH keys for host image injection.
1. [Host creation](./host-creation/README.md): Create one or more hosts.
1. [Volume creation](./volume-creation/README.md): Create iSCSI volumes for host attachments.
1. [Network creation](./network-creation/README.md): Create custom new metworks for project intra-communication.
1. [Usage information](./usage/README.md): Extract resource usage information.
1. [Project](./project/README.md): Create and manipulate projects.
1. [Advanced](./advanced/README.md): Advanced Terraform operations.
