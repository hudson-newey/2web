package builders

import (
	"github.com/hudson-newey/2web-cli/src/packages"
	"github.com/hudson-newey/2web-cli/src/shell"
)

func ServeSolution() {
	packageManager := packages.DeterminePackageManager()

	switch packageManager {
	case packages.Npm:
		npmServe()
	case packages.Pnpm:
		pnpmServe()
	case packages.Yarn:
		yarnServe()
	case packages.Bun:
		bunServe()
	}
}

func bunServe() {
	shell.ExecuteCommand(
		"bun",
		"exec",
		"vite",
		".",
		"--config",
		viteConfigLocation(),
	)
}

func npmServe() {
	shell.ExecuteCommand(
		"npm",
		"exec",
		"vite",
		".",
		"--config",
		viteConfigLocation(),
	)
}

func pnpmServe() {
	shell.ExecuteCommand(
		"pnpm",
		"exec",
		"vite",
		".",
		"--config",
		viteConfigLocation(),
	)
}

func yarnServe() {
	shell.ExecuteCommand(
		"yarn",
		"exec",
		"vite",
		".",
		"--config",
		viteConfigLocation(),
	)
}
