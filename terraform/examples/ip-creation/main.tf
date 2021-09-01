// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

provider "quake" {

}

resource "quake_ip" "ip" {
  ip_pool_id = var.ip_pool_id
  ip         = var.ip
  usage      = "Usage for ip"
}
