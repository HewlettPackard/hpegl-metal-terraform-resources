
// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

output "pnet" {
  # Output the created network
  value = hpegl_metal_network.pnet
}

output "ip" {
  # Output the allocated IP
  value = hpegl_metal_ip.ip
}
