
// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

output "pnet" {
  # Output the created network
  value = quake_network.pnet
}

output "ip" {
  # Output the allocated IP
  value = quake_ip.ip
}
