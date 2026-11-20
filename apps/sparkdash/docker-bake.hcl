target "docker-metadata-action" {}

variable "APP" {
  default = "sparkdash"
}

// Upstream publishes no releases and stopped tagging at v2.32.0, so VERSION
// mirrors package.json at REVISION and is bumped by hand with it.
variable "VERSION" {
  default = "2.70.0"
}

variable "REVISION" {
  // renovate: datasource=git-refs depName=https://github.com/MiaAI-Lab/sparkDash branch=main
  default = "2c736a7261ea2a589f256dc1b2fe133201cf0dce"
}

variable "SOURCE" {
  default = "https://github.com/MiaAI-Lab/sparkDash"
}

group "default" {
  targets = ["image-local"]
}

target "image" {
  inherits = ["docker-metadata-action"]
  args = {
    VERSION  = "${VERSION}"
    REVISION = "${REVISION}"
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
