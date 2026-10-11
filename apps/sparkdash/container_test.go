package main

import (
	"testing"

	"github.com/kalenarndt/containers/testhelpers"
)

func Test(t *testing.T) {
	image := testhelpers.GetTestImage("ghcr.io/kalenarndt/sparkdash:rolling")
	testhelpers.TestHTTPEndpoint(t, image, testhelpers.HTTPTestConfig{
		Port:       "5555",
		Path:       "/api/health",
		StatusCode: 200,
	}, nil)
}
