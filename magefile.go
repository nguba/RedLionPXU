//go:build mage
// +build mage

package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime" // To handle potential .exe suffix on Windows

	"github.com/magefile/mage/mg" // Mage utilities
	"github.com/magefile/mage/sh" // Shell helper for common commands
)

// Proto generates Go code from .proto files.
func Proto() error {
	fmt.Println("Generating Go code from .proto files...")

	// Use Go's native os.MkdirAll to create the directory, completely bypassing shell mkdir issues.
	if err := os.MkdirAll("public/api/pxu/v1", 0755); err != nil {
		return fmt.Errorf("failed to create proto output directory: %w", err)
	}

	// Define the protoc command and its arguments.
	// sh.RunV is a Mage helper that runs a command with verbose output.
	// It handles cross-platform path separators and executable finding.
	return sh.RunV("protoc",
		"--proto_path=internal/proto",
		"--go_out=public/api/pxu/v1",
		"--go_opt=paths=source_relative",
		"--go-grpc_out=public/api/pxu/v1",
		"--go-grpc_opt=paths=source_relative",
		"internal/proto/pxu.proto",
	)
}

// Build compiles the server and client binaries.
func Build() error {
	mg.Deps(Proto) // Declare dependency: Build needs Proto to run first

	fmt.Println("Building pxu-grpc-server...")

	if err := sh.RunV("go", "build", "-o", buildOutputPath("pxu-grpc-server"), "./cmd/pxu-grpc-server"); err != nil {
		return err
	}

	fmt.Println("Building pxu-cli...")
	return sh.RunV("go", "build", "-o", buildOutputPath("pxu-cli"), "./cmd/pxu-cli")
}

// Run builds and runs the gRPC server.
func Run() error {
	mg.Deps(Build) // Declare dependency: Run needs Build to run first

	fmt.Println("Running pxu-grpc-server...")
	// Use exec.Command for running the built binary directly, capturing output.
	cmd := exec.Command(buildOutputPath("pxu-grpc-server"))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Test runs all Go tests.
func Test() error {
	fmt.Println("Running tests...")
	return sh.RunV("go", "test", "./...")
}

// Clean removes built binaries and generated code.
func Clean() error {
	fmt.Println("Cleaning up...")
	// Use Go's native os.RemoveAll to clean up directories, bypassing shell rm/rmdir issues.
	if err := os.RemoveAll("build"); err != nil {
		fmt.Printf("Warning: failed to remove build directory: %v\n", err)
	}
	if err := os.RemoveAll("public/api"); err != nil {
		fmt.Printf("Warning: failed to remove generated files: %v\n", err)
	}

	// For go clean, we can still use sh.RunV
	return sh.RunV("go", "clean", "-cache", "-modcache")
}

// Default target runs build.
var Default = Build

// Helper function to get correct build output path for Windows .exe
func buildOutputPath(name string) string {
	if runtime.GOOS == "windows" {
		return fmt.Sprintf("build/%s.exe", name)
	}
	return fmt.Sprintf("build/%s", name)
}
