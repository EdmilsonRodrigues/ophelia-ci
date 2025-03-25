package git

import (
        "embed"
        "fmt"
        "log"
        "os"
        "os/exec"
        "path/filepath"
)

//go:embed templates/*
var templates embed.FS

func CreateGitRepository(repoPath string) error {
        if err := createBareGitRepository(repoPath); err != nil {
                return fmt.Errorf("failed creating bare repo: %w", err)
        }

        templateContent, err := templates.ReadFile("templates/post-receive")
        if err != nil {
                return fmt.Errorf("failed to read post-receive template: %w", err)
        }

        if err := createPostReceiveHook(repoPath, templateContent); err != nil {
                return fmt.Errorf("failed creating post-receive hook: %w", err)
        }

        return nil
}

func createBareGitRepository(repoPath string) error {
        originalPath, err := os.Getwd()
        if err != nil {
                return fmt.Errorf("failed to get current working directory: %w", err)
        }

        if err := os.MkdirAll(repoPath, 0755); err != nil {
                return fmt.Errorf("failed to create repository directory: %w", err)
        }

        if err := os.Chdir(repoPath); err != nil {
                return fmt.Errorf("failed to change directory: %w", err)
        }

        cmd := exec.Command("git", "init", "--bare")
        output, err := cmd.CombinedOutput()
        if err != nil {
                return fmt.Errorf("failed to initialize bare repository: %w\n%s", err, output)
        }

        log.Printf("Bare Git repository %s created successfully!\n", repoPath)

        if err := os.Chdir(originalPath); err != nil {
                return fmt.Errorf("failed to change back to original directory: %w", err)
        }

        return nil
}

func createPostReceiveHook(repoPath string, templateContent []byte) error {
        hooksPath := filepath.Join(repoPath, "hooks")
        postReceivePath := filepath.Join(hooksPath, "post-receive")

        if err := os.MkdirAll(hooksPath, 0755); err != nil {
                return fmt.Errorf("failed to create hooks directory: %w", err)
        }

        if err := os.WriteFile(postReceivePath, templateContent, 0755); err != nil {
                return fmt.Errorf("failed to write post-receive file: %w", err)
        }

        fmt.Println("post-receive hook created successfully.")
        return nil
}
