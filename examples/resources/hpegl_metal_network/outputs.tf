// (C) Copyright 2020-2022, 2025 Hewlett Packard Enterprise Development LP

output "pnet" {
  # Output the created network.
  value = hpegl_metal_network.pnet
}

output "pnet1" {
  # Output the second created network.
  value = hpegl_metal_network.pnet1
}
