target "docker-metadata-action" {}

variable "APP" {
  default = "scrapling"
}

variable "VERSION" {
  // renovate: datasource=pypi depName=scrapling
  default = "0.4.8"
}

variable "SOURCE" {
  default = "https://github.com/D4Vinci/Scrapling"
}

group "default" {
  targets = ["image-local"]
}

target "image" {
  inherits = ["docker-metadata-action"]
  args = {
    VERSION = "${VERSION}"
  }
  labels = {
    "org.opencontainers.image.source" = "${SOURCE}"
  }
}

target "image-local" {
  inherits = ["image"]
  output = ["type=docker"]
  tags = ["${APP}:${VERSION}"]
}

target "image-all" {
  inherits = ["image"]
  platforms = [
    "linux/amd64"
  ]
}
