package cmd

import (
	"fmt"
	"project-manager/pkg"
	_ "strings"
)

func SyncCmd(config *pkg.Config) error {
	gitManager := pkg.NewGitManager(config)
	scanner := pkg.NewScanner(config, *gitManager)

	for projectName := range config.Projects {
		fmt.Printf("%-32s\tscan...\r", projectName)

		status, err := scanner.Scan(projectName)

		if err != nil {
			fmt.Printf("%-32s\terror: %s...\n", projectName, err.Error())
			continue
		}

		if status.IsCloned {
			if status.Uncommitted {
				err := gitManager.Reset(projectName)

				if err != nil {
					return err
				}
			}

			fmt.Printf("%-32s\t%-8s\t%-8s\r", projectName, status.CurrentBranch, "fetch...")
			err := gitManager.Fetch(projectName)

			if err != nil {
				return err
			}

			fmt.Printf("%-32s\t%-8s\t%-8s\r", projectName, status.CurrentBranch, "pull...")

			err = gitManager.Pull(projectName)

			if err != nil {
				return err
			}

			fmt.Printf("%-32s\t%-8s\t%-8s\n", projectName, status.CurrentBranch, "synced")
		} else {
			fmt.Printf("%-32s\t%-8s\t%-8s\r", projectName, status.CurrentBranch, "clone...")

			err := gitManager.CloneProject(projectName)

			if err != nil {
				return err
			}

			status, err := scanner.Scan(projectName)

			if err != nil {
				fmt.Printf("%-32s\terror: %s...\n", projectName, err.Error())
				continue
			}

			fmt.Printf("%-32s\t%-8s\t%-8s\n", projectName, status.CurrentBranch, "cloned")
		}
	}

	return nil
}
