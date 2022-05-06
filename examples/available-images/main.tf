// (C) Copyright 2020-2022 Hewlett Packard Enterprise Development LP

provider "hpegl" {
  metal {
    gl_token = false
  }
}

data "hpegl_metal_available_images" "centos" {
  filter {
    name   = "flavor"
    values = ["(?i)centos"] // case insensitive for Centos or centos etc.
  }
  filter {
    name   = "version"
    values = ["7.6.*"] // al 7.6.XXXX image variants
  }
}
