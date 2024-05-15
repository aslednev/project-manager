package cmd

import (
	"fmt"
	"project-manager/pkg"
	_ "strings"
)

func CheckoutCmd(branchName string, config *pkg.Config) error {
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

			displayBranch := status.CurrentBranch
			isChanged, err := gitManager.CheckoutBranch(projectName, branchName)

			if err != nil {
				fmt.Printf("%-32s\terror: %s...\n", projectName, err.Error())
				continue
			}

			if isChanged {
				displayBranch = branchName
			}

			fmt.Printf("%-32s\t%-8s\n", projectName, displayBranch)
		} else {
			fmt.Printf("%-32s\t%-8s\n", projectName, "no cloned")
		}
	}

	return nil
}
