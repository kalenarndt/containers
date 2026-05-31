package main

import (
	"testing"

	"github.com/kalenarndt/containers/testhelpers"
)

func Test(t *testing.T) {
	image := testhelpers.GetTestImage("ghcr.io/kalenarndt/tqm:rolling")
	testhelpers.TestCommandSucceeds(t, image, nil, "/app/bin/tqm", "version")
}
