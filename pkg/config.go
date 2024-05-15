package pkg

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	_ "io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Instruction struct {
	Build     []string `yaml:"build"`
	Provision []string `yaml:"provision"`
}

type Project struct {
	Repository struct {
		Source string `yaml:"source"`
		Branch string `yaml:"branch"`
		Target string `yaml:"target"`
	} `yaml:"repository"`

	Instructions Instruction `yaml:"instructions"`
}

type Config struct {
	FileName string
	Version  string
	Global   struct {
		Instructions Instruction         `yaml:"instructions"`
		Scripts      map[string][]string `yaml:"scripts"`
	} `yaml:"global"`
	Projects map[string]*Project `yaml:"projects"`
}

func LoadConfig(configPath *string) *Config {
	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		fmt.Println("Not enough arguments")
		os.Exit(1)
	}

	filePath := *configPath
	if stat, err := os.Stat(filePath); err == nil && stat.IsDir() {
		filePath = filepath.Join(filePath, "manager.yaml")
	}

	if !strings.HasSuffix(filePath, ".yaml") {
		fmt.Println("Config file must be a .yaml file")
		os.Exit(1)
	}

	conf := &Config{
		FileName: filePath,
	}

	if err := conf.Load(); err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Config file not found at path:", filePath)
		} else {
			fmt.Println("Error loading config:", err)
		}
		os.Exit(1)
	}

	return conf
}

func (c *Config) Load() error {
	configFile, err := os.Open(c.FileName)

	if err != nil {
		return err
	}

	defer func() {
		if _err := configFile.Close(); _err != nil && err == nil {
			err = _err
		}
	}()

	decoder := yaml.NewDecoder(configFile)
	return decoder.Decode(&c)
}

func (c *Config) GetProjectByName(projectName string) (*Project, error) {
	cfg, exists := c.Projects[projectName]

	if !exists {
		return nil, fmt.Errorf("project `%s` not present", projectName)
	}

	return cfg, nil
}
