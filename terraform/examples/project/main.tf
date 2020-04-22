provider "quake" {
  
}

resource "quake_project" "project" {
    name = "blob"
    profile = {
        company = "ACME"
        address = "Area51"
        email = "acme@intergalactic.universe"
        phone_number = "+112 234 1245 3245"
        project_description = "Primitive Life"
        project_name = "Umbrella Corporation"
    }
    limits = {
        hosts = 5
        volumes = 10
        volume_capacity = 300
    }
}
