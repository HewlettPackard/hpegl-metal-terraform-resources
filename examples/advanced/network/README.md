# Example of creating a shared network and allocate an IP

This is an example of creating a shared network and allocating an IP from the IP pool associated with the network.

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

  # hpegl_metal_ip.ip will be created
  + resource "hpegl_metal_ip" "ip" {
      + id         = (known after apply)
      + ip         = "10.0.0.4"
      + ip_pool_id = (known after apply)
      + usage      = "Usage for ip"
    }

  # hpegl_metal_network.pnet will be created
  + resource "hpegl_metal_network" "pnet" {
      + description = "A description of pnet"
      + host_use    = (known after apply)
      + id          = (known after apply)
      + ip_pool_id  = (known after apply)
      + kind        = (known after apply)
      + location    = "USA:Central:V2DCC01"
      + location_id = (known after apply)
      + name        = "pnet"

      + ip_pool {
          + base_ip     = "10.0.0.0"
          + description = "A description of npool"
          + dns         = []
          + ip_ver      = "IPv4"
          + name        = "npool"
          + netmask     = "/24"
          + ntp         = []
        }
    }

Plan: 2 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + ip    = {
      + id         = (known after apply)
      + ip         = "10.0.0.4"
      + ip_pool_id = (known after apply)
      + usage      = "Usage for ip"
    }
  + pnet = {
      + description = "A description of pnet"
      + host_use    = (known after apply)
      + id          = (known after apply)
      + ip_pool     = [
          + {
              + base_ip       = "10.0.0.0"
              + default_route = ""
              + description   = "A description of npool"
              + dns           = []
              + ip_ver        = "IPv4"
              + name          = "npool"
              + netmask       = "/24"
              + no_proxy      = ""
              + ntp           = []
              + proxy         = ""
              + sources       = []
            },
        ]
      + ip_pool_id  = (known after apply)
      + kind        = (known after apply)
      + location    = "USA:Central:V2DCC01"
      + location_id = (known after apply)
      + name        = "pnet"
    }

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

hpegl_metal_network.pnet: Creating...
hpegl_metal_network.pnet: Creation complete after 0s [id=9a784c00-1fab-41b2-bbad-9be6836494c9]
hpegl_metal_ip.ip: Creating...
hpegl_metal_ip.ip: Creation complete after 0s [id=0c2047ef-a24a-432a-b87d-7c454c1e3a83:10.0.0.4]

Apply complete! Resources: 2 added, 0 changed, 0 destroyed.

Outputs:

ip = {
  "id" = "0c2047ef-a24a-432a-b87d-7c454c1e3a83:10.0.0.4"
  "ip" = "10.0.0.4"
  "ip_pool_id" = "0c2047ef-a24a-432a-b87d-7c454c1e3a83"
  "usage" = "Usage for ip"
}
pnet = {
  "description" = "A description of pnet"
  "host_use" = "Optional"
  "id" = "9a784c00-1fab-41b2-bbad-9be6836494c9"
  "ip_pool" = toset([
    {
      "base_ip" = "10.0.0.0"
      "default_route" = ""
      "description" = "A description of npool"
      "dns" = tolist([])
      "ip_ver" = "IPv4"
      "name" = "npool"
      "netmask" = "/24"
      "no_proxy" = ""
      "ntp" = tolist([])
      "proxy" = ""
      "sources" = tolist([])
    },
  ])
  "ip_pool_id" = "0c2047ef-a24a-432a-b87d-7c454c1e3a83"
  "kind" = ""
  "location" = "USA:Central:V2DCC01"
  "location_id" = "1ad98170-993e-4bfc-8b84-e689ea9a429b"
  "name" = "pnet"
}

```

### Argument Reference

The following arguments are supported:

- `ip` - The IP address to be allocated.
- `ip_pool_id` - ID of the IP pool from which the IP will be alloacted.
- `usage` - Description of the IP allocation.