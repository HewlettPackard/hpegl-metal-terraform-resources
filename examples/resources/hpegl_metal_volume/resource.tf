# (C) Copyright 2020-2023 Hewlett Packard Enterprise Development LP

provider "hpegl" {
  metal {
    gl_token = false
  }
}

variable "location" {
  default = "USA:CO:FTC"
}

resource "hpegl_metal_volume" "test_vols" {
  count        = 1
  name         = "vol-${count.index}"
  size         = 20
  shareable    = true
  flavor       = "NVMe"
  location     = var.location
  volume_collection = "d5a63736-a03f-4779-8a08-0b3763f86704"
  description  = "Terraformed volume"
}
