// (C) Copyright 2020-2022 Hewlett Packard Enterprise Development LP

variable "start" {
  default = "2020-01-13T07:44:00Z"
}

variable "end" {
  default = "2020-04-13T07:44:00Z"
}

data "hpegl_metal_usage" "used" {
  start = var.start
  #end = var.end
}

output "compute_consumption" {
  value = data.hpegl_metal_usage.used.host_usage
}

output "volume_consumption" {
  value = data.hpegl_metal_usage.used.volume_usage
}
