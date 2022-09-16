terraform {
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
}