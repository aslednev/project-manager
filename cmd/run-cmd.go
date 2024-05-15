package cmd

import (
	"fmt"
	"project-manager/pkg"
	_ "strings"
)

func RunCmd(config *pkg.Config, logger *pkg.Logger, name string) error {
	builder := pkg.NewBuilderManager(config)

	scripts, exists := config.Global.Scripts[name]

	if !exists {
		return fmt.Errorf("script %s not found", name)
	}

	for _, step := range scripts {
		err := builder.RunCommand(step, "./")
		if err != nil {
			fmt.Printf("script %s error: %s...\n", name, err.Error())
			continue
		}
	}

	return nil
}
