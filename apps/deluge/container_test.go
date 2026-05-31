package main

import (
	"testing"

	"github.com/kalenarndt/containers/testhelpers"
)

func Test(t *testing.T) {
	image := testhelpers.GetTestImage("ghcr.io/kalenarndt/deluge:rolling")
	t.Run("HTTP endpoint test", func(t *testing.T) {
		testhelpers.TestHTTPEndpoint(t, image,
			testhelpers.HTTPTestConfig{
				Port: "8112",
			},
			&testhelpers.ContainerConfig{
				Env: map[string]string{
					"DELUGE_BIN": "deluge-web",
				},
			},
		)
	})
}
