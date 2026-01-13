package main

import (
	"os"
	"os/exec"
	"testing"
)

func TestIsGitRepo(t *testing.T) {
	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "gitx-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Change to temp directory
	oldDir, _ := os.Getwd()
	defer os.Chdir(oldDir)
	os.Chdir(tmpDir)

	// Should not be a git repo initially
	if isGitRepo() {
		t.Error("Expected not to be a git repo")
	}

	// Initialize git repo
	exec.Command("git", "init").Run()

	// Should be a git repo now
	if !isGitRepo() {
		t.Error("Expected to be a git repo after init")
	}
}

func TestGetGitConfig(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "gitx-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	oldDir, _ := os.Getwd()
	defer os.Chdir(oldDir)
	os.Chdir(tmpDir)

	exec.Command("git", "init").Run()
	exec.Command("git", "config", "user.name", "Test User").Run()

	name, err := getGitConfig("user.name")
	if err != nil {
		t.Fatalf("Failed to get git config: %v", err)
	}
	if name != "Test User" {
		t.Errorf("Expected 'Test User', got '%s'", name)
	}
}

