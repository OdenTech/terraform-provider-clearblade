variable "clearblade-creds" {
  type    = string
  default = "clearblade_service_account_json_auth_file"
}

variable "gcp_project_id" {
  type    = string
  default = "gcp_project_id"
}

variable "gcp_region" {
  type    = string
  default = "gcp_region_here"
}

variable "registry_id" {
  type    = string
  default = "registry_id_here"
}

variable "event_subfolder_matches" {
  type    = string
  default = "test/path"
}

variable "event_topic_name" {
  type    = string
  default = "projects/gcp_project_id_here/topics/rootevent"
}

variable "state_topic_name" {
  type    = string
  default = "projects/gcp_project_id_here/topics/rootevent"
}

variable "log_level" {
  type    = string
  default = "INFO"
}
