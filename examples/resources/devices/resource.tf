# https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/cloudiot_device

resource "google_cloudiot_device" "basic-device" {
  name     = "basic-device"
  registry = google_cloudiot_registry.iot-registry.id
}

# https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/cloudiot_device

resource "google_cloudiot_device" "advanced-device" {
  name     = "advanced-device"
  registry = google_cloudiot_registry.iot-registry.id

  credentials {
    public_key {
      format = "RSA_X509_PEM"
      key    = file("~/.auth/rsa_cert.pem")
    }
  }
}
