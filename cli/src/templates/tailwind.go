package templates

import (
	"github.com/hudson-newey/2web-cli/src/builders/configs"
	"github.com/hudson-newey/2web/_shared/logger"
)

func TailwindTemplate() {
	if !configs.HasViteConfig() {
		logger.PrintError(
			"2Web support for Tailwind is currently limited to Vite builds.\n" +
				"Please run '2web template vite' to add a Vite builder.",
		)
	}

	copyFromTemplates("tailwind.config.js", "tailwind.config.js")
	copyFromTemplates("postcss.config.js", "postcss.config.js")
}
