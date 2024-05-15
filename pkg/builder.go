package pkg

import (
	"os"
	"os/exec"
	"strings"
)

type BuilderManager struct {
	config *Config
}

func NewBuilderManager(config *Config) *BuilderManager {
	return &BuilderManager{config: config}
}

func (b *BuilderManager) Provision(projectName string) error {
	cfg, err := b.config.GetProjectByName(projectName)

	if err != nil {
		return err
	}

	for _, step := range cfg.Instructions.Provision {
		name := ""
		var args []string
		if strings.HasPrefix(step, "self") {
			name = os.Args[0]
			args = append(strings.Split(step, " ")[1:], projectName)
		} else {
			parts := strings.Fields(step)
			name = parts[0]
			args = parts[1:]
		}

		err = b.exec(name, args, cfg.Repository.Target)

		if err != nil {
			return err
		}
	}

	return nil
}

func (b *BuilderManager) StartBuildInstruction(projectName string) error {
	cfg, err := b.config.GetProjectByName(projectName)

	if err != nil {
		return err
	}

	for _, step := range cfg.Instructions.Build {
		err = b.RunCommand(step, cfg.Repository.Target)

		if err != nil {
			return err
		}
	}

	return nil
}

func (b *BuilderManager) RunCommand(command string, dir string) error {
	parts := strings.Fields(command)
	name := parts[0]
	args := parts[1:]

	return b.exec(name, args, dir)
}

func (b *BuilderManager) exec(name string, args []string, dir string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	return err
}
