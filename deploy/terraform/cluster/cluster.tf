resource "google_container_cluster" "cluster" {
  name = "antifreeze-cluster"
  zone = "${var.zone}"

  master_auth = {
    username = "${var.master_username}"
    password = "${var.master_password}"
  }

  node_pool = [
    {
      name       = "antifreeze-node-pool"
      node_count = 1

      node_config {
        image_type   = "COS"
        disk_size_gb = "10"
        machine_type = "g1-small"
        tags         = ["${var.back_end_tag}"]

        # Minimum required scopes for GKE VMs
        # Doesn't relate to container permissions
        # Actually, it's unclear whether these relate to container permissions
        oauth_scopes = [
          "compute-rw",
          "storage-ro",
          "logging-write",
          "monitoring",
          "datastore",
        ]
      }
    },
  ]

  addons_config {
    # The current architecture doesn't support horizontal scaling
    horizontal_pod_autoscaling {
      disabled = true
    }

    # A load balancer is used to keep a static IP
    http_load_balancing {
      disabled = false
    }

    # There hasn't been a need for the dashboard yet
    kubernetes_dashboard {
      disabled = true
    }

    # Hasn't been a need for it yet
    network_policy_config {
      disabled = true
    }
  }
}
