package main

import (
	"testing"

	"github.com/kalenarndt/containers/testhelpers"
)

func Test(t *testing.T) {
	image := testhelpers.GetTestImage("ghcr.io/kalenarndt/opentofu-runner:rolling")
	testhelpers.TestFileExists(t, image, "/usr/local/bin/terraform", nil)
}
