// (C) Copyright 2021-2022 Hewlett Packard Enterprise Development LP

provider "hpegl" {
  metal {
    gl_token = false
  }
}

resource "hpegl_metal_ip" "ip" {
  ip_pool_id = var.ip_pool_id
  ip         = var.ip
  usage      = "Usage for ip"
}
