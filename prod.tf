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
		username  = "${var.master_username}"
		password  = "${var.master_password}"
	}

  node_config {
		image_type   = "COS"
		disk_size_gb = "10"
    machine_type = "g1-small"
		tags         = ["back-end"]

    oauth_scopes = [
      "compute-rw",
      "storage-ro",
      "logging-write",
      "monitoring",
    ]
  }
}
