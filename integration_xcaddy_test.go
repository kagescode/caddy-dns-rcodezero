package rcodezero

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

func TestDockerAdaptBuild(t *testing.T) {
	if os.Getenv("CADDY_DNS_DOCKER_INTEGRATION") == "" {
		t.Skip("set CADDY_DNS_DOCKER_INTEGRATION=1 to run docker integration test")
	}

	if _, err := exec.LookPath("docker"); err != nil {
		t.Skip("docker not found in PATH")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	image := "caddy-rcodezero-adapt-test:local"

	// Build using your Dockerfile.adapt-test (which runs 'caddy adapt' during build).
	cmd := exec.CommandContext(ctx, "docker", "build",
		"-f", "Dockerfile.adapt-test",
		"-t", image,
		".",
	)

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	t.Logf("Running: %s", strings.Join(cmd.Args, " "))
	if err := cmd.Run(); err != nil {
		t.Fatalf("docker build failed: %v\n%s", err, out.String())
	}

	t.Logf("docker build ok:\n%s", out.String())
}

func TestDockerRunSmoke(t *testing.T) {
	if os.Getenv("CADDY_DNS_DOCKER_INTEGRATION_RUN") == "" {
		t.Skip("set CADDY_DNS_DOCKER_INTEGRATION_RUN=1 to run docker run smoke test")
	}
	if os.Getenv("CADDY_DNS_DOCKER_INTEGRATION") == "" {
		t.Skip("set CADDY_DNS_DOCKER_INTEGRATION=1 as well (build step required)")
	}

	if _, err := exec.LookPath("docker"); err != nil {
		t.Skip("docker not found in PATH")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	image := "caddy-rcodezero-adapt-test:local"

	// Ensure image exists (build it). Reuse the same docker build logic.
	buildCmd := exec.CommandContext(ctx, "docker", "build",
		"-f", "Dockerfile.adapt-test",
		"-t", image,
		".",
	)
	var buildOut bytes.Buffer
	buildCmd.Stdout = &buildOut
	buildCmd.Stderr = &buildOut
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("docker build failed: %v\n%s", err, buildOut.String())
	}

	// Run the container briefly (it should exit quickly because this Dockerfile
	// doesn't define a CMD; but we can still execute a harmless command to prove the binary runs).
	// We just run: caddy version
	runCtx, runCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer runCancel()

	runCmd := exec.CommandContext(runCtx, "docker", "run", "--rm", image, "caddy", "version")
	var runOut bytes.Buffer
	runCmd.Stdout = &runOut
	runCmd.Stderr = &runOut

	t.Logf("Running: %s", strings.Join(runCmd.Args, " "))
	if err := runCmd.Run(); err != nil {
		t.Fatalf("docker run failed: %v\n%s", err, runOut.String())
	}

	t.Logf("docker run ok:\n%s", runOut.String())
}

