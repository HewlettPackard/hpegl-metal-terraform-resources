provider "quake" {

}

resource "quake_host" "terra_host" {
  count         = 1
  name          = "tformed-${count.index}"
  image_flavor  = "centos"                
  image_version = "7.6.1810"  
  machine_size  = "Any"
  ssh           = ["User1 - Linux"]  
  networks      = ["Private", "Public", "Storage"]  
  location      = var.location
  description   = "Hello from Terraform"
  # This will create and attach an iSCSI volume to the host.
  volumes {
    name   = "large-volume-${count.index}"
    size   = 5
    flavor = "Fast"
  }
  # Attach additional volumes by using multiple volume{} blocks.
  #volumes {
  #  name   = "small-volume-${count.index}"
  #  size   = 2
  #  flavor = "Fast"
  #}
}
