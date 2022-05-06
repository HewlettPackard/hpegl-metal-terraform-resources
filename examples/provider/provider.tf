# (C) Copyright 2022 Hewlett Packard Enterprise Development LP

# Set-up for terraform >= v0.13
terraform {
  required_providers {
    hpegl = {
      source  = "terraform.example.com/metal/hpegl"
      version = ">= 0.0.1"
    }
  }
}

# Example of provider configuration when using GreenLake IAM token
provider "hpegl" {
  metal {
    rest_url   = "https://localhost:3002"
    space_name = "space_name"
    project_id = "1d96bfbc-9cf0-4268-aac6-ca1c65aca385"
  }
}

# Example of provider configuration when using Metal Service token
provider "hpegl" {
  metal {
    rest_url = "https://localhost:3002"
    gl_token = false
  }
}
