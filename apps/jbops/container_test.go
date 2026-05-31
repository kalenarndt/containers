package main

import (
	"testing"

	"github.com/kalenarndt/containers/testhelpers"
)

func Test(t *testing.T) {
	image := testhelpers.GetTestImage("ghcr.io/kalenarndt/jbops:rolling")
	testhelpers.TestFileExists(t, image, "/app/fun/plexapi_haiku.py", nil)
}
