# hpegl-metal-terraform-resources

- [hpegl-metal-terraform-resources](#hpegl-metal-terraform-resources)
- [Terraform resources for HPEGL Metal](#terraform-resources-for-hpegl-metal)
  - [Requirements](#requirements)
  - [Usage](#usage)
  - [Building the resources as stand-alone provider](#building-the-resources-as-stand-alone-provider)

# Terraform resources for HPEGL Metal

Terraform Metal resources is a plugin for HPEGL terraform provider that allows the full lifecycle management of HPEGL
Metal resources. This provider is maintained by [HPEGL Metal resources team](mailTo:quake-core@hpe.com).

## Requirements

1. Terraform version >= v0.13 [install terraform](https://learn.hashicorp.com/tutorials/terraform/install-cli)
2. A Service Client to authenticate against GreenLake.
3. Terraform basics. [Terraform Introduction](https://www.terraform.io/intro/index.html)

## Usage

See the terraform provider for
hpegl [documentation](https://registry.terraform.io/providers/HewlettPackard/hpegl/latest/docs)
to get started using the provider.

## Building the resources as stand-alone provider

```bash
$ git clone https://github.com/hewlettpackard/hpegl-metal-terraform-resources.git
$ cd hpegl-metal-terraform-resources
$ make build
```

Note: For debugging the provider please refer to the
[debugging guide](https://medium.com/@gandharva666/debugging-terraform-using-jetbrains-goland-f9a7e992cb1d)
