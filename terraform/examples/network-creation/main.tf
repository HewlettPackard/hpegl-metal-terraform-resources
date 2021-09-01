// (C) Copyright 2016-2021 Hewlett Packard Enterprise Development LP

provider "quake" {

}

resource "quake_network" "pnet" {
    name = "pnet"
    description = "A description of pnet"
    location = var.location
     ip_pool {
      name="npool"
      description="A description of npool"
      ip_ver= "IPv4"
      base_ip= "10.0.0.0"
      netmask= "/24"
      default_route = "10.0.0.1"
      sources {
        base_ip="10.0.0.3"
        count = 10
      }
      dns = [ "10.0.0.50" ]
      proxy = "10.0.0.60"
      no_proxy = "10.0.0.5"
      ntp = [ "10.0.0.80" ]
    }
}
