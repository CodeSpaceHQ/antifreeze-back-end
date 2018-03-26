variable master_password {
  type        = "string"
  description = "Kubernetes master password"
}

variable master_username {
  type        = "string"
  description = "Kubernetes master username"
}

variable target_tags {
  default = "back-end"
}
