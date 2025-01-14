package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Tillter2998/newGo/internal/applicationStrategy"
)

func main() {
	var projectType string
	var projectName string
	var projectDir string
	var githubUser string

	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	flag.StringVar(&projectType, "type", "empty", "Type of project to bootstrap")
	flag.StringVar(&projectName, "name", "newGoApp", "Name of the project")
	flag.StringVar(&projectDir, "dir", currentDir, "Directory to create project in")
	flag.StringVar(&githubUser, "githubUser", "", "Github user to initialize go packages for")

	flag.Parse()

	appContext := new(applicationStrategy.ApplicationContext)
	appRegistry := applicationStrategy.GetRegistry()

	if flag.NFlag() > 0 {
		// run without TUI

		// Slice setup to easily expand required flags later on
		missingFlags := make([]string, 0, 1)

		if githubUser == "" {
			missingFlags = append(missingFlags, "githubUser")
		}

		if len(missingFlags) > 0 {
			log.Fatal(fmt.Sprintf("Missing the following required flags: %s", strings.Join(missingFlags, ", ")))
		}

		strategy, err := appRegistry.GetStrategy(projectType)
		if err != nil {
			log.Fatal(err)
		}

		appContext.SetStrategy(strategy)
		err = appContext.CreateApplication(githubUser, projectName, projectDir)
		if err != nil {
			log.Fatal(err)
		}
		return
	} else {
		// TODO: Add TUI element
	}
}
