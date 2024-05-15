package cmd

import (
	"fmt"
	"project-manager/pkg"
	_ "strings"
)

func ScanCmd(config *pkg.Config) error {
	gitManager := pkg.NewGitManager(config)
	scanner := pkg.NewScanner(config, *gitManager)

	for projectName := range config.Projects {
		fmt.Printf("%-32s\tscan...\r", projectName)

		status, err := scanner.Scan(projectName)

		if err != nil {
			fmt.Printf("%-32s\terror: %s...\n", projectName, err.Error())
			continue
		}

		hasChanged := ""

		if status.Uncommitted {
			hasChanged = "has local changes"
		}

		if status.IsCloned {
			fmt.Printf("%-32s\t%-8s\t%-12s\n", projectName, status.CurrentBranch, hasChanged)
		} else {
			fmt.Printf("%-32s\t%-8s\n", projectName, "no cloned")
		}
	}

	return nil
}
