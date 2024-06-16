package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	ProjectsHome string
	Projects     []string
}

func readProjectsHome() string {
	file, err := os.Open(".prjrc")
	if err != nil {
		fmt.Println("Error opening .prjrc file:", err)
		os.Exit(1)
	}
	defer file.Close()

	var projectsHome string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "PROJECTS_HOME=") {
			projectsHome = strings.TrimPrefix(line, "PROJECTS_HOME=")
			projectsHome = strings.TrimSpace(projectsHome)
			projectsHome = os.ExpandEnv(projectsHome)
			break
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading .prjrc file:", err)
		os.Exit(1)
	}

	return projectsHome
}

func GetConfig() Config {
	projectsHome := readProjectsHome()

	projects, err := filepath.Glob(projectsHome)

	if err != nil {
		fmt.Println("Error getting projects:", err)
		os.Exit(1)
	}

	return Config{
		ProjectsHome: projectsHome,
		Projects:     projects,
	}
}
