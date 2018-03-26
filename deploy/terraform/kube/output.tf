output "instance_group_urls" {
  value = "${google_container_cluster.c.instance_group_urls}"
}

output "endpoint" {
	value = "${google_container_cluster.c.endpoint}"
}

output "client_certificate" {
	value = "${google_container_cluster.c.master_auth.0.client_certificate}"
}

output "client_key" {
	value = "${google_container_cluster.c.master_auth.0.client_key}"
}

output "cluster_ca_certificate" {
	value = "${google_container_cluster.c.master_auth.0.cluster_ca_certificate}"
}
