terraform {
  required_providers {
    inpyu = {
      source = "registry.terraform.io/study/inpyu-ossca"
    }
  }
}

provider "inpyu" {
  username = "inpyu"
  password = "test123"
  host     = "http://localhost:19090"
}


resource "inpyu_cafe" "cafe" {
    name = "Sample Cafe"
    address = "123 Coffee St"
    description = "A cozy place to enjoy coffee and pastries"
    image = "http://example.com/image.jpg"
}