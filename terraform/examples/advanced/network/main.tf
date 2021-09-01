// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

provider "quake" {

}

resource "quake_network" "pnet" {
  name        = "pnet"
  description = "A description of pnet"
  location    = var.location
  ip_pool {
    name        = "npool"
    description = "A description of npool"
    ip_ver      = "IPv4"
    base_ip     = "10.0.0.0"
    netmask     = "/24"
  }
}

resource "quake_ip" "ip" {
  ip_pool_id = quake_network.pnet.ip_pool_id
  ip         = "10.0.0.4"
  usage      = "Usage for ip"
}
