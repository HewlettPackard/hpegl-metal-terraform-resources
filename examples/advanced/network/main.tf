// (C) Copyright 2021-2022 Hewlett Packard Enterprise Development LP

provider "hpegl" {
  metal {
    gl_token = false
  }
}

resource "hpegl_metal_network" "pnet" {
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

resource "hpegl_metal_ip" "ip" {
  ip_pool_id = hpegl_metal_network.pnet.ip_pool_id
  ip         = "10.0.0.4"
  usage      = "Usage for ip"
}
