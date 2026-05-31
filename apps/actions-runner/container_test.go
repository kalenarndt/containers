package main

import (
	"testing"

	"github.com/kalenarndt/containers/testhelpers"
)

func Test(t *testing.T) {
	image := testhelpers.GetTestImage("ghcr.io/kalenarndt/actions-runner:rolling")
	testhelpers.TestFileExists(t, image, "/usr/local/bin/yq", nil)
}
