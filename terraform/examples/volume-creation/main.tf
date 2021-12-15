# (C) Copyright 2020-2021 Hewlett Packard Enterprise Development LP

provider "quake" {

}

resource "quake_volume" "test_vols" {
  count       = 1
  name        = "vol-${count.index}"
  size        = 20
  shareable   = true
  flavor      = "Fast"
  location    = var.location
  description = "Terraformed volume"
}
