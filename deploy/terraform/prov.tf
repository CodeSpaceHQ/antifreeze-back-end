# TODO(NilsG-S): Set versions for tooling

provider "google" {
  # Uses application default credentials
  alias   = "default"
  project = "antifreeze-199016"
  region  = "us-central1"
  zone    = "us-central1-a"
}

# TODO: Aliasing causes problems with implicit stuff (see addr)
provider "google" {
  project = "antifreeze-199016"
  region  = "us-central1"
  zone    = "us-central1-a"
}

provider "kubernetes" {
  host     = "${module.kube.endpoint}"
  username = "${var.master_username}"
  password = "${var.master_password}"

  client_certificate     = "${base64decode(module.kube.client_certificate)}"
  client_key             = "${base64decode(module.kube.client_key)}"
  cluster_ca_certificate = "${base64decode(module.kube.cluster_ca_certificate)}"
}
