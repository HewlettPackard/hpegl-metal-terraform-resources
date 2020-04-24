provider "quake" {

}

resource "quake_network" "pnet" {
    name = "pnet"
    description = "A description of pnet"
    location = var.location
}
