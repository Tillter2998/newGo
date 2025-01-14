package main

import (
	"flag"
	"log"
	"os"

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
	}
}
