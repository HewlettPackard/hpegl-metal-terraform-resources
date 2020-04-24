# Example of creating a network

This is an example of creating a network for use by a project.

To run the example:
* Authenticate against a portal using steeld login
* Update `variables.tf` OR provide overrides on the command line
* Run with a command similar to
```
terraform apply \
    -var ="location=USA:Texas:AUSL2"
``` 

## Example output

```
An execution plan has been generated and is shown below.
Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # quake_network.pnet will be created
  + resource "quake_network" "pnet" {
      + description = "A desciption of pnet"
      + host_use    = (known after apply)
      + id          = (known after apply)
      + kind        = (known after apply)
      + location    = "USA:Texas:AUSL2"
      + location_id = (known after apply)
      + name        = "pnet"
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

quake_network.pnet: Creating...
quake_network.pnet: Creation complete after 0s [id=edaa63c1-3f01-47a8-933f-cbee5a30708f]

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

Outputs:

pnet = {
  "description" = "A desciption of pnet"
  "host_use" = "Optional"
  "id" = "edaa63c1-3f01-47a8-933f-cbee5a30708f"
  "kind" = "Custom"
  "location" = "USA:Texas:AUSL2"
  "name" = "pnet"
}

```

### Argument Reference

The following arguments are supported:

- `name` - The name of the network.
- `description` - (Optional) Some descriptive text that helps describe the network and purpose.
- `location` - Where the network is to be created in country:region:data-center style.


### Attribute Reference

In addition to the arguments listed above, the following attributes are exported:

- `location_id` - Unique ID of the location.
- `kind` - The kind of network, e.g. "Custom".
- `host_use` - The requirement of a host to use this network, e.g. "Required" or "Optional"




