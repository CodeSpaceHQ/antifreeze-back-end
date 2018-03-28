# TODO(NilsG-S): Set versions for tooling

provider "google" {
  # Uses application default credentials
  version = "~> 1.8"
  project = "${var.project_id}"
  region  = "${var.region}"
  zone    = "${var.zone}"
}

provider "google" {
  version = "~> 1.8"
  alias   = "cluster"
  project = "${var.project_id}"
  region  = "${var.region}"
  zone    = "${var.zone}"
}

provider "kubernetes" {
  version  = "~> 1.1"
  host     = "${module.cluster.endpoint}"
  username = "${var.master_username}"
  password = "${var.master_password}"

  client_certificate     = "${base64decode(module.cluster.client_certificate)}"
  client_key             = "${base64decode(module.cluster.client_key)}"
  cluster_ca_certificate = "${base64decode(module.cluster.cluster_ca_certificate)}"
}

# Apparently backend configs can't contain interpolation
terraform {
  backend "gcs" {
    bucket  = "antifreeze-tf-state"
    prefix  = "prod"
    project = "antifreeze-199016"
    region  = "us-central1"
  }
}
