package cmd

import (
	"fmt"
	"project-manager/pkg"
	_ "strings"
)

func BuildCmd(config *pkg.Config, logger *pkg.Logger, name *string) error {
	gitManager := pkg.NewGitManager(config)
	scanner := pkg.NewScanner(config, *gitManager)
	builder := pkg.NewBuilderManager(config)

	var processedProjects = make(map[string]string)

	for projectName := range config.Projects {
		if name != nil && projectName != *name {
			continue
		}

		fmt.Printf("%-32s\tscan...\r", projectName)

		status, err := scanner.Scan(projectName)

		if err != nil {
			fmt.Errorf("%-32s\terror: %s...\n", projectName, err.Error())
			processedProjects[projectName] = err.Error()
			continue
		}

		if !status.IsCloned {
			fmt.Printf("%-32s\t%-8s\t%-8s\r", projectName, "-", "not cloned")
			processedProjects[projectName] = "not cloned"
			continue
		}

		fmt.Printf("%-32s\t%-8s\t%-8s\n", projectName, status.CurrentBranch, "building...")

		err = builder.StartBuildInstruction(projectName)

		if err != nil {
			processedProjects[projectName] = err.Error()
			logger.Fatal(err)
		} else {
			processedProjects[projectName] = "Build on " + status.CurrentBranch
		}
	}

	fmt.Println("\nProcessed projects:")
	for projectName := range processedProjects {
		fmt.Printf("%-32s\t%-32s\n", projectName, processedProjects[projectName])
	}

	return nil
}
