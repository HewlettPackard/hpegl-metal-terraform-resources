output "compute_consumption" {
  value = data.hpegl_metal_usage.used.host_usage
}

output "volume_consumption" {
  value = data.hpegl_metal_usage.used.volume_usage
}
