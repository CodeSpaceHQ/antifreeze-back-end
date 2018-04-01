variable project_id {
  type        = "string"
  description = "GCP project ID"
  default     = "antifreeze-199016"
}

variable region {
  type        = "string"
  description = "GCP project region"
  default     = "us-central1"
}

variable zone {
  type        = "string"
  description = "GCP project zone"
  default     = "us-central1-a"
}

variable back_end_tag {
  default = "back-end"
}

variable master_password {
  type        = "string"
  description = "Kubernetes master password"
}

variable master_username {
  type        = "string"
  description = "Kubernetes master username"
}

variable antifreeze_secret {
  type        = "string"
  description = "Signing secret for JWT"
}
