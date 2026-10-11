package main

import (
	"testing"

	helpers "github.com/kalenarndt/containers/tests"
)

func Test(t *testing.T) {
	image := helpers.GetTestImage("ghcr.io/kalenarndt/sparkdash:rolling")
	helpers.RequireHTTPEndpoint(t, image, helpers.HTTPTestConfig{
		Port:       "5555",
		Path:       "/api/health",
		StatusCode: 200,
	}, nil)
}
