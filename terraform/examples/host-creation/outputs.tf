output "ips" {
  # Output a map of hostame with all the IP addresses assigned on each network.
  value = zipmap(quake_host.terra_host.*.name,  quake_host.terra_host.*.connections)
}
