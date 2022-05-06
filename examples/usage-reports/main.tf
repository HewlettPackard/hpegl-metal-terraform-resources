// (C) Copyright 2020-2022 Hewlett Packard Enterprise Development LP

provider "hpegl" {
  metal {
    gl_token = false
  }
}
data "hpegl_metal_usage" "used" {
  start = var.start
  #end = var.end
}
