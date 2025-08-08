// (C) Copyright 2020-2022, 2025 Hewlett Packard Enterprise Development LP

provider "hpegl" {
  metal {
    gl_token = false
  }
}

variable "location" {
  // Provide a location at which to query for resources. 
  // The format is country:region:data-center
  default = "USA:Central:AFCDCC1"
}

resource "hpegl_metal_network" "pnet" {
  name        = "TestPrivateNet001"
  description = "Created by TF"
  host_use    = "Default"
  location    = var.location
  ip_pool {
    name          = "PrivNet1Pool"
    description   = "IPPool for TestPrivateNet001"
    ip_ver        = "IPv4"
    base_ip       = "10.35.0.0"
    netmask       = "/24"
    default_route = "10.35.0.1"
    sources {
      base_ip = "10.35.0.3"
      count   = 10
    }
    dns      = ["8.8.8.8"]
    proxy    = "http://10.0.0.60:8080"
    no_proxy = "10.35.0.5"
    ntp      = ["10.35.0.80"]
  }
  vlan = 4000
  vni  = 40400
}

resource "hpegl_metal_network" "pnet1" {
  name        = "TestPrivateNet002"
  description = "Created by TF - no IP Pool"
  purpose     = "Backup"
  location    = var.location
  no_ip_pool  = true
}
