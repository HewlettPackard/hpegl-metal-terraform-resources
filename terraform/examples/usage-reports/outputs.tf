output  "compute_consumption" {
    value = data.quake_usage.used.host_usage
}

output  "volume_consumption" {
    value = data.quake_usage.used.volume_usage
}