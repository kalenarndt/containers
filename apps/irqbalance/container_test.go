package main

import (
	"testing"

	"github.com/kalenarndt/containers/testhelpers"
)

func Test(t *testing.T) {
	image := testhelpers.GetTestImage("ghcr.io/kalenarndt/irqbalance:rolling")
	testhelpers.TestFileExists(t, image, "/usr/sbin/irqbalance", nil)
}
