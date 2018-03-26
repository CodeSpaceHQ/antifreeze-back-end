# TODO(NilsG-S): Set versions for tooling

provider "google" {
  # Uses application default credentials
  project = "${var.project_id}"
  region  = "${var.region}"
  zone    = "${var.zone}"
}

provider "google" {
  alias   = "cluster"
  project = "${var.project_id}"
  region  = "${var.region}"
  zone    = "${var.zone}"
}

provider "kubernetes" {
  host     = "${module.kube.endpoint}"
  username = "${var.master_username}"
  password = "${var.master_password}"

  client_certificate     = "${base64decode(module.kube.client_certificate)}"
  client_key             = "${base64decode(module.kube.client_key)}"
  cluster_ca_certificate = "${base64decode(module.kube.cluster_ca_certificate)}"
}
