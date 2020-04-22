provider "quake" {
}

data "quake_available_images" "centos" {
  filter {
       name = "flavor"
       values = ["(?i)centos"]    // case insensitive for Centos or centos etc.
  }
  filter {
      name = "version"
      values = ["7.6.*"]  // al 7.6.XXXX image variants
  }
}