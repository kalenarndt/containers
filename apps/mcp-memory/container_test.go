package main

import (
	"testing"

	"github.com/kalenarndt/containers/testhelpers"
)

func Test(t *testing.T) {
	image := testhelpers.GetTestImage("ghcr.io/kalenarndt/mcp-memory:rolling")

	t.Run("Check mcp-server-memory exists", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, image, nil, "/usr/local/bin/mcp-server-memory", "--version")
	})
}
