package main

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

func TestHelloWorldBinary(t *testing.T) {
	binary := os.Getenv("HELLO_BINARY")
	if binary == "" {
		binary = filepath.Join(t.TempDir(), "hello-world")
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()
		build := exec.CommandContext(ctx, "go", "build", "-o", binary, ".")
		build.Env = append(os.Environ(), "CGO_ENABLED=0")
		if output, err := build.CombinedOutput(); err != nil {
			t.Fatalf("build failed: %v\n%s", err, output)
		}
	}
	binary, err := filepath.Abs(binary)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, binary)
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("binary failed: %v\nstderr: %s", err, stderr.String())
	}
	if got := stdout.String(); got != "Hello World\n" {
		t.Errorf("stdout = %q; want %q", got, "Hello World\n")
	}
	if got := stderr.String(); got != "" {
		t.Errorf("stderr = %q; want empty output", got)
	}
}
