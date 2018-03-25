provider "google" {
  # Uses application default credentials
  project = "antifreeze-199016"
  zone    = "us-central1-a"
}

resource "google_container_cluster" "c" {
  name               = "antifreeze-c"
  zone               = "us-central1-a"
  initial_node_count = 1

  node_config {
    machine_type = "g1-small"

    oauth_scopes = [
      "compute-rw",
      "storage-ro",
      "logging-write",
      "monitoring",
    ]
  }
}
