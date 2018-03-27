module "cluster" {
  source = "./cluster"

  master_username = "${var.master_username}"
  master_password = "${var.master_password}"
  back_end_tag    = "${var.back_end_tag}"
  zone            = "${var.zone}"

  providers = {
    google = "google.cluster"
  }
}

# TODO: Add secrets to container
# resource "kubernetes_secret" "ksec" {}

# Service account for app containers
resource "google_service_account" "container" {
  account_id   = "container"
  display_name = "container"
}

resource "google_service_account_key" "container-key" {
  service_account_id = "${google_service_account.container.id}"
}

/*
	google_service_account_iam_* resources are for controlling
	access to service accounts as a resource. google_project_iam
	resources are for controlling roles as they apply to a
	service account.
*/
resource "google_project_iam_member" "container-iam" {
  role   = "roles/datastore.user"
  member = "serviceAccount:${google_service_account.container.email}"
}

# TODO: put kubernetes stuff in module?
resource "kubernetes_service" "service" {
  metadata {
    name = "antifreeze-service"
  }

  spec {
    selector {
      App = "${kubernetes_pod.pod.metadata.0.labels.App}"
    }

    port {
      port        = 8081
      target_port = 8081
    }

    type = "LoadBalancer"

    # Assign static IP to this service's load balancer
    load_balancer_ip = "${google_compute_address.addr.address}"
  }
}

resource "kubernetes_pod" "pod" {
  metadata {
    name = "antifreeze-pod"

    # Used to select this pod in kubernetes_service
    labels {
      App = "antifreeze"
    }
  }

  spec {
    container {
      # Make sure to keep this updated!
      image = "nilsgs/antifreeze:1"
      name  = "antifreeze-container"

      # List of ports to expose
      port {
        # This is the port for the server
        container_port = 8081
      }
    }
  }
}

# Configuration of static IP
resource "google_compute_address" "addr" {
  name = "antifreeze-addr"

  # This specifies that the address is static and external
  address_type = "EXTERNAL"
}
