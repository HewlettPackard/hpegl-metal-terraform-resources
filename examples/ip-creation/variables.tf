// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

variable "location" {
  // Provide a location at which to query for resources. The default given here
  // is compatible with the portal-simulator.
  // the format is country:region:data-center
  default = "USA:Texas:AUSL2"
}

variable "ip_pool_id" {
  // Provide id of the IP pool form which the IP will be allocated.
}

variable "ip" {
  // Provide the IP to be allocated from the IP pool.
}
