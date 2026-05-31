package main

import (
	"testing"

	"github.com/kalenarndt/containers/testhelpers"
)

func Test(t *testing.T) {
	image := testhelpers.GetTestImage("ghcr.io/kalenarndt/busybox:rolling")
	testhelpers.TestCommandSucceeds(t, image, nil, "/bin/busybox", "--list")
}
