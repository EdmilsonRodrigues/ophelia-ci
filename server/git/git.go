package git

import (
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

//go:embed templates/*
var templates embed.FS

// CreateGitRepository initializes a new bare Git repository at the specified path
// and sets up a post-receive hook.
//
// It performs the following steps:
// 	1. Creates a bare Git repository in the given directory.
// 	2. Reads the post-receive template content from the embedded templates.
// 	3. Creates a post-receive hook using the template content.
//
// If any step fails, an error is returned with details.
func CreateGitRepository(repoPath, gitignore string) error {
	if err := createBareGitRepository(repoPath); err != nil {
		return fmt.Errorf("failed creating bare repo: %w", err)
	}

	if err := runInitialCommit(repoPath, gitignore); err != nil {
		return fmt.Errorf("failed running initial commit: %w", err)
	}

	if err := createPostReceiveHook(repoPath); err != nil {
		return fmt.Errorf("failed creating post-receive hook: %w", err)
	}

	return nil
}

// UpdateGitRepository updates an existing Git repository by renaming its path.
//
// The function will:
//   - Rename the directory at the given repository path to the new path.
//
// If the renaming fails, an error is returned with details.
func UpdateGitRepository(repoPath, newPath string) error {
	err := os.Rename(repoPath, newPath)
	if err != nil {
		return fmt.Errorf("failed to rename repository: %w", err)
	}
	return nil
}

// DeleteGitRepository removes a Git repository at the given path.
//
// The function will:
//   - Remove the directory at the given path and all its contents
//
// If the removal fails, an error is returned with details.
func DeleteGitRepository(repoPath string) error {
	err := os.RemoveAll(repoPath)
	if err != nil {
		return fmt.Errorf("failed to remove repository: %w", err)
	}
	return nil
}

// createBareGitRepository creates a bare Git repository at the given path.
//
// The function will:
// - Change directory to the given path
// - Initialize a bare Git repository
// - Change back to the original directory
//
// If any of the above steps fail, an error is returned.
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

// createPostReceiveHook creates a post-receive hook file in the hooks directory
// of the given repository path. The file content is the given templateContent.
//
// The function will:
//   - Create the hooks directory if it does not exist
//   - Write the content of the template to a post-receive file in the hooks
//     directory
//
// If any of the above steps fail, an error is returned.
func createPostReceiveHook(repoPath string) error {
	hooksPath := filepath.Join(repoPath, "hooks")
	postReceivePath := filepath.Join(hooksPath, "post-receive")

	templateContent, err := templates.ReadFile("templates/post-receive")
	if err != nil {
		return fmt.Errorf("failed to read post-receive template: %w", err)
	}

	if err := os.MkdirAll(hooksPath, 0755); err != nil {
		return fmt.Errorf("failed to create hooks directory: %w", err)
	}

	if err := os.WriteFile(postReceivePath, templateContent, 0755); err != nil {
		return fmt.Errorf("failed to write post-receive file: %w", err)
	}

	fmt.Println("post-receive hook created successfully.")
	return nil
}

// runInitialCommit creates a temporary directory and initializes a regular
// git repository inside it. It then creates a .gitignore file with the given
// gitignore and commits it with the message "Initial commit". Finally, it pushes
// the commits to the given repository path.
//
// If any of the above steps fail, an error is returned with details.
func runInitialCommit(repoPath, gitignore string) error {
	originalPath, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %w", err)
	}
	defer os.Chdir(originalPath)

	tempDir, err := createTempInitedDir(repoPath)
	if err != nil {
		return fmt.Errorf("failed to create temporary directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	gitignorePath, err := createBaseGitignore(tempDir, gitignore)
	if err != nil {
		return fmt.Errorf("failed creating .gitignore: %w", err)
	}

	if err := commitFile(tempDir, "Initial commit", gitignorePath); err != nil {
		return fmt.Errorf("failed to commit .gitignore: %w", err)
	}

	if err := pushCommitsToRemote(tempDir, repoPath); err != nil {
		return fmt.Errorf("failed to push commits to remote: %w", err)
	}

	return nil
}

// pushCommitsToRemote adds a remote repository and pushes commits to it.
//
// This function performs the following steps:
// 1. Saves the current working directory.
// 2. Changes the working directory to the specified repository path.
// 3. Adds a remote named "origin" with the given remote path.
// 4. Pushes commits to the master branch of the remote repository.
//
// If any of these steps fail, an error is returned with details.
func pushCommitsToRemote(repoPath, remotePath string) error {
	originalPath, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %w", err)
	}
	defer os.Chdir(originalPath)

	if err := os.Chdir(repoPath); err != nil {
		return fmt.Errorf("failed to change directory: %w", err)
	}
	remoteCmd := exec.Command("git", "remote", "add", "origin", remotePath)
	remoteOutput, err := remoteCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to add remote: %w\n%s", err, remoteOutput)
	}
	pushCmd := exec.Command("git", "push", "-u", "origin", "master")
	pushOutput, err := pushCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to push commits: %w\n%s", err, pushOutput)
	}
	return nil
}

// createTempInitedDir creates a temporary directory and initializes a regular
// git repository inside it. The temporary directory is created in the system's
// default temporary directory and has the given repository path as its base
// name. The function then changes into the temporary directory and runs the
// command "git init" to initialize a regular git repository. If any of the
// above steps fail, an error is returned.
func createTempInitedDir(repoPath string) (string, error) {
	originalPath, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current working directory: %w", err)
	}
	defer os.Chdir(originalPath)

	tempDir, err := os.MkdirTemp("", "git-repo")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary directory: %w", err)
	}

	if err := os.Chdir(tempDir); err != nil {
		return "", fmt.Errorf("failed to change directory: %w", err)
	}

	initCmd := exec.Command("git", "init")
	initOutput, err := initCmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to initialize regular git repository in %s: %w\n%s", repoPath, err, initOutput)
	}

	return tempDir, nil

}

// createBaseGitignore creates a .gitignore file in the given repository path
// with content based on the given stack.
//
// The function will:
//   - Fetch the .gitignore template content for the given stack from
//     https://github.com/github/gitignore/blob/main/<stack>.gitignore
//   - Extract the JSON payload from the HTML response
//   - Unmarshal the JSON into the GigtignorePage struct
//   - Write the raw lines of the payload to a .gitignore file in the given
//     repository path
//
// If any of the above steps fail, an error is returned.
func createBaseGitignore(repoPath string, stack string) (string, error) {
	gitignoreURL := fmt.Sprintf("https://github.com/github/gitignore/blob/main/%s.gitignore", stack)
	gitignorePath := filepath.Join(repoPath, ".gitignore")

	resp, err := http.Get(gitignoreURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch gitignore template for %s: %w", stack, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch gitignore template for %s, status code: %d", stack, resp.StatusCode)
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read gitignore template for %s: %w", stack, err)
	}
	bodyString := string(content)

	re := regexp.MustCompile(`<script type="application/json" data-target="react-app\.embeddedData">(.*?)</script>`)
	match := re.FindStringSubmatch(bodyString)

	if len(match) < 2 {
		return "", fmt.Errorf("script tag with data-target 'react-app.embeddedData' not found")
	}

	scriptContent := match[1]

	var data GigtignorePage
	err = json.Unmarshal([]byte(scriptContent), &data)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	gitignoreContent := strings.Join(data.Payload.Blob.RawLines, "\n")

	err = os.WriteFile(gitignorePath, []byte(gitignoreContent), 0644)
	if err != nil {
		return "", fmt.Errorf("failed to write .gitignore file: %w", err)
	}

	log.Printf(".gitignore created for %s stack in repo %s.\n", stack, repoPath)
	return gitignorePath, nil
}

// commitFile commits a file to the repository at repoPath with a commit message.
//
// The function will:
// - Change directory to the given repository path
// - Run `git add <filePath>`
// - Run `git commit -m "<message>"`
//
// If any of the above steps fail, an error is returned.
func commitFile(repoPath, message, filePath string) error {
	originalPath, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %w", err)
	}
	defer os.Chdir(originalPath)

	if err := os.Chdir(repoPath); err != nil {
		return fmt.Errorf("failed to change directory to repo: %w", err)
	}

	addCmd := exec.Command("git", "add", filePath)
	addOutput, err := addCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to run git add %s: %w\n%s", filePath, err, addOutput)
	}

	commitCmd := exec.Command("git", "commit", "-m", message)
	commitOutput, err := commitCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to run git commit -m \"%s\": %w\n%s", message, err, commitOutput)
	}

	log.Printf("Committed %s to repo %s with message: %s\n", filePath, repoPath, message)
	return nil
}

type GigtignorePage struct {
	Payload struct {
		Blob struct {
			RawLines []string `json:"rawLines"`
		} `json:"blob"`
	} `json:"payload"`
}
