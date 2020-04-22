provider "quake" {

}

resource "quake_volume" "test_vols" {
  count       = 1
  name        = "vol-${count.index}"
  size        = 20
  flavor      = "Fast" 
  location    = var.location
  description = "Terraformed volume"
}
