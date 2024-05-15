package pkg

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"log"
	"os"
	"os/exec"
	"strings"
)

type GitManager struct {
	config *Config
}

func NewGitManager(config *Config) *GitManager {
	return &GitManager{config: config}
}

func (g *GitManager) CloneProject(projectName string) error {

	homeDir, err := os.UserHomeDir()

	if err != nil {
		return err
	}

	var sshKeyPath = homeDir + "/.ssh/id_rsa"
	sshAuth, err := ssh.NewPublicKeysFromFile("git", sshKeyPath, "")
	if err != nil {
		log.Fatal(err)
	}

	cfg, err := g.config.GetProjectByName(projectName)

	if err != nil {
		return err
	}

	_, err = git.PlainClone(cfg.Repository.Target, false, &git.CloneOptions{
		URL:  cfg.Repository.Source,
		Auth: sshAuth,
	})

	if err != nil {
		return err
	}

	return nil
}

func (g *GitManager) Fetch(projectName string) error {
	cfg, err := g.config.GetProjectByName(projectName)

	if err != nil {
		return err
	}

	cmd := exec.Command("git", "fetch")
	cmd.Dir = cfg.Repository.Target
	_, err = cmd.Output()

	return err
}

func (g *GitManager) CheckoutBranch(projectName string, branchName string) (bool, error) {
	cfg, err := g.config.GetProjectByName(projectName)

	if err != nil {
		return false, err
	}

	cmd := exec.Command("git", "checkout", branchName)
	cmd.Dir = cfg.Repository.Target

	_, err = cmd.Output()

	isChanged := true

	if err != nil {
		isChanged = false
	}

	return isChanged, nil
}

func (g *GitManager) Reset(projectName string) error {
	cfg, err := g.config.GetProjectByName(projectName)

	if err != nil {
		return err
	}

	cmd := exec.Command("git", "checkout", "--", ".")
	cmd.Dir = cfg.Repository.Target
	_, err = cmd.Output()

	return err
}

func (g *GitManager) HasUncommitedChanges(projectName string) (bool, error) {
	cfg, err := g.config.GetProjectByName(projectName)

	if err != nil {
		return false, err
	}

	repository, err := git.PlainOpen(cfg.Repository.Target)

	if err != nil {
		return false, err
	}

	worktree, err := repository.Worktree()

	if err != nil {
		return false, fmt.Errorf("cannot get workdir: %w", err)
	}

	status, err := worktree.Status()

	if err != nil {
		return false, fmt.Errorf("cannot check status: %w", err)
	}

	noChanges := true

	if !status.IsClean() {
		lines := strings.Split(status.String(), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)

			if line == "" || strings.HasPrefix(line, "?") {
				continue
			}

			noChanges = false
		}
	}

	return !noChanges, nil
}

func (g *GitManager) GetCurrentBranch(projectName string) (string, error) {
	cfg, err := g.config.GetProjectByName(projectName)

	if err != nil {
		return "", err
	}

	//git rev-parse --abbrev-ref HEAD

	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = cfg.Repository.Target
	output, err := cmd.Output()

	if err != nil {
		return "", fmt.Errorf("cannot get HEAD: %w", err)
	}

	branchName := strings.TrimSpace(string(output))

	return branchName, nil
}

func (g *GitManager) Pull(projectName string) error {
	cfg, err := g.config.GetProjectByName(projectName)

	if err != nil {
		return err
	}

	cmd := exec.Command("git", "pull")
	cmd.Dir = cfg.Repository.Target
	_, err = cmd.Output()

	return err
}
