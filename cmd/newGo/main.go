package main

import (
	"flag"
	"log"

	"github.com/Tillter2998/newGo/internal/applicationStrategy"
)

func main() {
	var projectType string
	var projectName string
	var projectDir string

	flag.StringVar(&projectType, "type", "", "Type of project to bootstrap")
	flag.StringVar(&projectName, "name", "", "Name of the project")
	flag.StringVar(&projectDir, "dir", "", "Directory to create project in")

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
		appContext.CreateApplication(projectName, projectDir)
		return
	}
}
