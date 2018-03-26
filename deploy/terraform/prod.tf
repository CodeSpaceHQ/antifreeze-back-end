module "p" {
  source = "./pool.tf"
}

provider "google" {
  # Uses application default credentials
  project = "antifreeze-199016"
  zone    = "us-central1-a"
}

resource "google_container_cluster" "c" {
  name               = "antifreeze-c"
  zone               = "us-central1-a"
  initial_node_count = 1

  master_auth = {
    username = "${var.master_username}"
    password = "${var.master_password}"
  }

  node_config {
    image_type   = "COS"
    disk_size_gb = "10"
    machine_type = "g1-small"
    tags         = ["${var.target_tags}"]

    oauth_scopes = [
      "compute-rw",
      "storage-ro",
      "logging-write",
      "monitoring",
    ]
  }

  addons_config {
    horizontal_pod_autoscaling {
      disabled = true
    }

    http_load_balancing {
      disabled = false
    }

    kubernetes_dashboard {
      disabled = false
    }

    network_policy_config {
      disabled = false
    }
  }
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
        group = "${module.p.np}"
      },
    ]
  }
}
