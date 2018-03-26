variable target_tags {
	default = "back-end"
}

resource "google_container_node_pool" "np" {
  name       = "antifreeze-np"
  zone       = "us-central1-a"
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
}

output "np" {
	value = "${google_container_node_pool.np}"
}
