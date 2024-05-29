terraform {
  required_providers {
    inpyu = {
      source = "registry.terraform.io/study/inpyu-ossca"
    }
  }
}


provider "inpyu" {
#  username = "inpyu1234"
#  password = "test123"
#  host     = "http://localhost:19090"
}

#resource "cafe" "cafe1" {
#    name = "Sample Cafe"
#    address = "123 Coffee St"
#    description = "A cozy place to enjoy coffee and pastries"
#    image = "http://example.com/image.jpg"
#}
