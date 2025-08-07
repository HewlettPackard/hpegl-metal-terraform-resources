# (C) Copyright 2020-2023, 2025 Hewlett Packard Enterprise Development LP

provider "hpegl" {
  metal {
    gl_token = false
  }
}

variable "location" {
  default = "USA:Central:AFCDCC1"
}

resource "hpegl_metal_volume" "test_vols" {
  count             = 1
  name              = "vol-${count.index}"
  size              = 20
  shareable         = true
  flavor            = "Block - Standard"
  location          = var.location
  description       = "Terraformed volume"
}
