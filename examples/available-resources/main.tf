provider "hpegl" {
  metal {
    gl_token = false
  }
}

data "hpegl_metal_available_resources" "physical" {

}
