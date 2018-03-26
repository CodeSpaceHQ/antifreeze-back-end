provider "google" {
  # Uses application default credentials
  project = "antifreeze-199016"
  zone    = "us-central1-a"
}

module "gce-lb-http" {
  source = "GoogleCloudPlatform/lb-http/google"
  name   = "antifreeze-lb"
  ssl    = false

  target_tags = ["${var.target_tags}"]

  url_map        = "${google_compute_url_map.um.self_link}"
  create_url_map = false

  backends = {
    "0" = [
      {
        group = "${module.comp.np}"
      },
    ]
  }
}
