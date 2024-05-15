package main

import (
	"flag"
	"fmt"
	"os"
	"project-manager/cmd"
	"project-manager/pkg"
)

func main() {
	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		fmt.Println("Not enough arguments")
		os.Exit(1)
	}

	configPath := flag.String("c", "manager.yaml", "Path to the configuration file or directory")

	conf := pkg.LoadConfig(configPath)

	log := pkg.NewLogger(pkg.WARNING)

	switch args[0] {
	case "scan":
		err := cmd.ScanCmd(conf)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		break
	case "checkout":
		if len(args) < 2 {
			fmt.Println("No branch name provided")
			os.Exit(1)
		}

		err := cmd.CheckoutCmd(args[1], conf)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		break
	case "sync":
		err := cmd.SyncCmd(conf)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		break
	case "build":
		var projectName *string = nil

		if len(args) > 1 {
			projectName = &args[1]
		}

		err := cmd.BuildCmd(conf, log, projectName)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		break
	case "provision":
		var projectName *string = nil

		if len(args) > 1 {
			projectName = &args[1]
		}

		err := cmd.ProvisionCmd(conf, log, projectName)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		break
	case "run":

		if len(args) < 2 {
			fmt.Println("No script name provided")
			os.Exit(1)
		}

		err := cmd.RunCmd(conf, log, args[1])

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		break
	default:
		fmt.Println("Invalid command")
		os.Exit(1)
	}

	os.Exit(0)
}
