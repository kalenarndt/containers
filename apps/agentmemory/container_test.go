package main

import (
	"context"
	"testing"

	"github.com/kalenarndt/containers/testhelpers"
)

func Test(t *testing.T) {
	ctx := context.Background()
	image := testhelpers.GetTestImage("ghcr.io/kalenarndt/agentmemory:rolling")
	testhelpers.TestHTTPEndpoint(t, ctx, image, testhelpers.HTTPTestConfig{
		Port:       "3111",
		Path:       "/agentmemory/livez",
		StatusCode: 200,
	}, nil)
}