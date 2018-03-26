module "kube" {
  source = "./kube"

  master_username = "${var.master_username}"
  master_password = "${var.master_password}"
  target_tags     = "${var.target_tags}"

  providers = {
    google = "google.default"
  }
}

# TODO: This may be unnecessary for kubernetes
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
        group = "${module.kube.instance_group_urls[0]}"
      },
    ]
  }

  backend_params = [
    "/,http,8081,10",
  ]
}

# TODO: This may be unnecessary for kubernetes
resource "google_compute_url_map" "um" {
  name            = "antifreeze-um"
  default_service = "${module.gce-lb-http.backend_services[0]}"

  host_rule = {
    hosts        = ["*"]
    path_matcher = "allpaths"
  }

  path_matcher = {
    name            = "allpaths"
    default_service = "${module.gce-lb-http.backend_services[0]}"

    path_rule {
      paths   = ["/test/post"]
      service = "${kubernetes_service.kser.metadata.self_link}"
    }
  }
}

# resource "kubernetes_secret" "ksec" {}

resource "kubernetes_service" "kser" {
  metadata {
    name = "antifreeze-kser"
  }

  spec {
    selector {
      App = "${kubernetes_pod.kp.metadata.0.labels.App}"
    }

    port {
      port        = 8081
      target_port = 8081
    }

    type = "LoadBalancer"
  }
}

resource "kubernetes_pod" "kp" {
  metadata {
    name = "antifreeze-kp"

    # Used to select this pod in kubernetes_service
    labels {
      App = "antifreeze"
    }
  }

  spec {
    container {
      image = "nilsgs/antifreeze"

      # Ensures that the container is updated
      image_pull_policy = "Always"
      name              = "antifreeze-kc"

      # List of ports to expose
      port {
        # This is the port for the server
        container_port = 8081
      }
    }
  }
}

# Configuration of static ip
# resource "google_compute_address" "addr" {
# 	name = "antifreeze-a"
# }

