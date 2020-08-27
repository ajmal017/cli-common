package client

import (
	"context"
	"testing"
)

func TestNewClientWithOpts(t *testing.T) {
	cli, err := NewClientWithOpts(WithHost("http://localhost:8080"))
	if err != nil {
		t.Errorf("%s", err)
	}

	ping, err := cli.Ping(context.Background())
	if err != nil {
		t.Fatalf("failed")
	}
	t.Errorf(ping.APIVersion)
}