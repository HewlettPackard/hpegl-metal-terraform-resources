provider "quake" {

}

data "quake_usage" "used" {
  start = var.start
  #end = var.end
}
