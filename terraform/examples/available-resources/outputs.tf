
output "infrastructure" {
  value = "${data.quake_available_resources.physical}"
}

# output "locations" {
#     value = data.quake_available_resources.physical.locations
# }

# output "images" {
#     value = data.quake_available_resources.physical.images
# }

# output "ssh-keys" {
#     value = data.quake_available_resources.physical.ssh_keys
# }

# output "networks" {
#     value = data.quake_available_resources.physical.networks #[for net in data.quake_available_resources.physical.networks : net if net.location == var.location]
# }

# output "volumes" {
#     value = [for vol in data.quake_available_resources.physical.volumes : vol if vol.location == var.location]
# }

# output "volume-flavors" {
#     value = data.quake_available_resources.physical.volume_flavors
# }

# output "machine-sizes" {
#     value = data.quake_available_resources.physical.machine_sizes 
# }