module "comp" {
	source = "./comp"
}

resource "google_container_cluster" "c" {
  name = "antifreeze-c"
  zone = "us-central1-a"

  master_auth = {
    username = "${var.master_username}"
    password = "${var.master_password}"
  }

  node_pool = [
		"${module.comp.np}"
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
