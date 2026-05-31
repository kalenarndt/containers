package main

import (
	"testing"

	"github.com/kalenarndt/containers/testhelpers"
)

func Test(t *testing.T) {
	image := testhelpers.GetTestImage("ghcr.io/kalenarndt/k8s-sidecar:rolling")
	testhelpers.TestCommandSucceeds(t, image, nil, "python", "-u", "-m", "sidecar", "--help")
}
