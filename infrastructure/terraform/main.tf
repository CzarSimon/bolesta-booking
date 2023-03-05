terraform {
  backend "s3" {
    endpoint                    = "https://s3.fr-par.scw.cloud"
    key                         = "terraform.tfstate"
    bucket                      = "bolesta-booking-infrastructure-terraform-state"
    region                      = "fr-par"
    skip_credentials_validation = true
    skip_region_validation      = true
  }
  required_providers {
    scaleway = {
      source = "scaleway/scaleway"
    }
  }
  required_version = ">= 0.13"
}

provider "scaleway" {
  zone       = "fr-par-1"
  region     = "fr-par"
  project_id = "ccdc8190-9ce3-427f-ba23-98000872dfec"
}

resource "scaleway_object_bucket" "frontend" {
  name = "bolesta-booking-frontend"
  acl  = "public-read"
}

resource "scaleway_object_bucket_website_configuration" "frontend_website" {
  bucket = scaleway_object_bucket.frontend.name
  index_document {
    suffix = "index.html"
  }
  error_document {
    key = "index.html"
  }
}

resource "scaleway_instance_ip" "server_ip" {}

resource "scaleway_domain_record" "booking_A_record" {
  dns_zone = "xn--blesta-wxa.se"
  name     = "booking"
  type     = "A"
  data     = scaleway_instance_ip.server_ip.address
  ttl      = 300
}

resource "scaleway_instance_security_group" "server_sg" {
  inbound_default_policy  = "drop"
  outbound_default_policy = "accept"

  inbound_rule {
    action   = "accept"
    protocol = "TCP"
    port     = "22"
  }

  inbound_rule {
    action   = "accept"
    protocol = "TCP"
    port     = "80"
  }

  inbound_rule {
    action   = "accept"
    protocol = "TCP"
    port     = "443"
  }
}

resource "scaleway_instance_server" "server" {
  name              = "main-server"
  type              = "DEV1-S"
  image             = "ubuntu_focal"
  ip_id             = scaleway_instance_ip.server_ip.id
  security_group_id = scaleway_instance_security_group.server_sg.id

  user_data = {
    cloud-init = file("${path.module}/cloud-init.yml")
  }
}
