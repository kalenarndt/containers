package main

import (
	"testing"

	"github.com/kalenarndt/containers/testhelpers"
)

func Test(t *testing.T) {
	image := testhelpers.GetTestImage("ghcr.io/kalenarndt/cni-plugins:rolling")
	testhelpers.TestCommandSucceeds(t, image, nil, "/plugins/macvlan")
}
