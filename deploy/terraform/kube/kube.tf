resource "google_container_cluster" "c" {
  name = "antifreeze-c"
  zone = "us-central1-a"

  master_auth = {
    username = "${var.master_username}"
    password = "${var.master_password}"
  }

  node_pool = [
    {
      name       = "antifreeze-np"
      node_count = 1

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
    },
  ]

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
