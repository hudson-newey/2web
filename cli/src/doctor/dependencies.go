package doctor

import (
	"fmt"
	"os"
	"os/exec"
	"text/tabwriter"
)

// Automatically checks that all dependencies are installed on the system.
func CheckDependencies() {
	requiredDependencies := []dependency{
		{name: "pandoc", url: "https://pandoc.org/installing.html"},
		{name: "sass", url: "https://sass-lang.com/install"},
		{name: "ffmpeg", url: "https://ffmpeg.org/download.html"},
		{name: "node", url: "https://nodejs.org/en/download"},
		{name: "npm", url: "https://docs.npmjs.com/cli/configuring-npm/install"},
		{name: "2webc", url: "https://github.com/hudson-newey/2web"},
	}

	tabWriter := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer tabWriter.Flush()

	tableHeader := "Dependency\tInstalled\n"
	_, err := tabWriter.Write([]byte(tableHeader))
	if err != nil {
		panic(err)
	}

	tabWriter.Write([]byte("----------\t---------\n"))

	for _, dep := range requiredDependencies {
		isInstalled := isDependencyInstalled(dep)
		status := dependencyStatusString(dep, isInstalled)

		_, err := tabWriter.Write([]byte(status))
		if err != nil {
			panic(err)
		}
	}
}

func isDependencyInstalled(dep dependency) bool {
	_, err := exec.LookPath(dep.name)
	return err == nil
}

func dependencyStatusString(dep dependency, isInstalled bool) string {
	const failedIcon string = "\033[31m✘\033[0m"
	const successIcon string = "\033[32m✔\033[0m"

	if isInstalled {
		return fmt.Sprintf("%s\t%s\n", dep.name, successIcon)
	}

	return fmt.Sprintf("%s\t%s\n", dep.name, failedIcon)
}
