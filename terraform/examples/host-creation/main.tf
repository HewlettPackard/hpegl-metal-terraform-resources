# (C) Copyright 2020-2021 Hewlett Packard Enterprise Development LP

provider "quake" {

}

resource "quake_host" "terra_host" {
  count         = 1
  name          = "tformed-${count.index}"
  image         = "centos@7.6.1810"
  machine_size  = "Any"
  ssh           = ["User1 - Linux"]
  networks      = ["Private", "Public", "Storage"]
  network_route = "Public"
  location      = var.location
  description   = "Hello from Terraform"
  # This will create a shareable iSCSI volume and attach it to the host.
  volumes {
    name        = "large-volume-${count.index}"
    size        = 5
    shareable   = true
    flavor      = "Fast"
    location    = var.location
    description = "Terraformed volume"
  }
  # Create and attach additional volumes by using multiple volume{} blocks.
  #volumes {
  #  name   = "small-volume-${count.index}"
  #  size   = 2
  #  flavor = "Fast"
  #}
}
