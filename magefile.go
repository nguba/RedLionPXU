//go:build mage
// +build mage

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

const (
	protoDir  = "public/api/pxu/v1"
	protoFile = "pxu.proto"
	buildDir  = "build"
	serverCmd = "./cmd/pxu-grpc-server"
	clientCmd = "./cmd/pxu-cli"
	serverBin = "pxu-grpc-server"
	clientBin = "pxu-cli"
)

// Proto generates Go code from .proto files.
func Proto() error {
	fmt.Println("Generating Go code from .proto files...")

	if err := os.MkdirAll(protoDir, 0755); err != nil {
		return fmt.Errorf("failed to create proto output directory: %w", err)
	}

	return sh.RunV("protoc",
		"--proto_path="+protoDir,
		"--go_out="+protoDir,
		"--go_opt=paths=source_relative",
		"--go-grpc_out="+protoDir,
		"--go-grpc_opt=paths=source_relative",
		filepath.Join(protoDir, protoFile),
	)
}

// Build compiles the server and client binaries.
func Build() error {
	mg.Deps(Proto)

	if err := os.MkdirAll(buildDir, 0755); err != nil {
		return fmt.Errorf("failed to create build directory: %w", err)
	}

	//if err := BuildServer(); err != nil {
	//	return err
	//}

	//if err := BuildClient(); err != nil {
	//	return err
	//}

	return nil
}

func BuildServer() error {
	fmt.Println("Building", serverBin+"...")
	if err := sh.RunV("go", "build", "-o", buildOutputPath(serverBin), serverCmd); err != nil {
		return fmt.Errorf("failed to build server: %w", err)
	}

	return nil
}

func BuildClient() error {
	fmt.Println("Building", clientBin+"...")
	if err := sh.RunV("go", "build", "-o", buildOutputPath(clientBin), clientCmd); err != nil {
		return fmt.Errorf("failed to build client: %w", err)
	}
	return nil
}

// Run builds and runs the gRPC server.
func Run() error {
	mg.Deps(Build)

	fmt.Println("Running", serverBin+"...")
	cmd := exec.Command(buildOutputPath(serverBin))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Test runs all Go tests.
func Test() error {
	mg.Deps(Build)
	fmt.Println("Running tests...")
	return sh.RunV("go", "test", "-v", "./...")
}

// TestCoverage runs tests with coverage report.
func TestCoverage() error {
	fmt.Println("Running tests with coverage...")
	return sh.RunV("go", "test", "-v", "-coverprofile=coverage.out", "./...")
}

// Lint runs golangci-lint on the project.
func Lint() error {
	fmt.Println("Running linter...")
	return sh.RunV("golangci-lint", "run")
}

// Format formats Go code using gofmt.
func Format() error {
	fmt.Println("Formatting code...")
	return sh.RunV("gofmt", "-s", "-w", ".")
}

// Clean removes built binaries and generated code.
func Clean() error {
	fmt.Println("Cleaning up...")

	// Remove build directory
	if err := os.RemoveAll(buildDir); err != nil {
		fmt.Printf("Warning: failed to remove build directory: %v\n", err)
	}

	// Clean generated proto files
	if err := CleanV1Api(); err != nil {
		return fmt.Errorf("failed to clean v1 API: %w", err)
	}

	// Remove coverage files
	removeFiles("coverage.out")

	// Clean Go cache
	return sh.RunV("go", "clean", "-cache")
}

// CleanV1Api removes generated Go files from proto compilation.
func CleanV1Api() error {
	fmt.Println("Cleaning up v1 API...")
	return removeFiles(filepath.Join(protoDir, "*.go"))
}

// Check runs all quality checks (format, lint, test).
func Check() error {
	mg.Deps(Format, Lint, Test)
	return nil
}

// Install installs the binaries to $GOPATH/bin.
func Install() error {
	mg.Deps(Build)

	fmt.Println("Installing binaries...")

	if err := sh.RunV("go", "install", serverCmd); err != nil {
		return fmt.Errorf("failed to install server: %w", err)
	}

	return sh.RunV("go", "install", clientCmd)
}

// Dev runs the development workflow: clean, build, test.
func Dev() error {
	mg.Deps(Clean, Build, Test)
	return nil
}

// Default target runs build.
var Default = Build

// Helper function to get correct build output path for Windows .exe
func buildOutputPath(name string) string {
	path := filepath.Join(buildDir, name)
	if runtime.GOOS == "windows" {
		return path + ".exe"
	}
	return path
}

// Helper function to remove files matching a pattern
func removeFiles(pattern string) error {
	files, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	for _, file := range files {
		if err := os.Remove(file); err != nil {
			// Don't fail if file doesn't exist
			if !os.IsNotExist(err) {
				return fmt.Errorf("failed to remove %s: %w", file, err)
			}
		}
	}

	return nil
}
