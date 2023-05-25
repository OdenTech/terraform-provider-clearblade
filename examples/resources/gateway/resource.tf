# https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/cloudiot_device

resource "google_cloudiot_device" "iot-gateway" {
  name     = "iot-gateway"
  registry = google_cloudiot_registry.iot-registry.id

  credentials {
    public_key {
      format = "RSA_X509_PEM"
      key    = file("~/.auth/rsa_cert.pem")
    }
  }

  gateway_config {
    gateway_type        = "GATEWAY"
    gateway_auth_method = "ASSOCIATION_ONLY"
  }
}
