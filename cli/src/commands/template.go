package commands

import (
	"fmt"

	"github.com/hudson-newey/2web-cli/src/templates"
	"github.com/hudson-newey/2web/_shared/logger"
)

func templateCommand(programName string, command string, args []string) {
	argsLen := len(args)
	if argsLen < 3 {
		errorMsg := fmt.Sprintf("invalid arguments:\n\texpected: %s %s <template>", programName, command)
		logger.PrintError(errorMsg)
	}

	template := args[2]

	templates := map[string]func(){
		"serer-side-rendering": templates.SsrTemplate,
		"ssr":                  templates.SsrTemplate,

		"database": templates.DatabaseTemplate,
		"db":       templates.DatabaseTemplate,

		"load-balancer": templates.LoadBalancerTemplate,
		"lb":            templates.LoadBalancerTemplate,

		"sitemap":     templates.SitemapTemplate,
		"sitemap.xml": templates.SitemapTemplate,

		"robots.txt": templates.RobotsTxtTemplate,

		"security.txt": templates.SecurityTemplate,

		"llms.txt": templates.LlmsTemplate,

		"vite":     templates.ViteTemplate,
		"tailwind": templates.TailwindTemplate,
	}

	templateFunc, exists := templates[template]
	if !exists {
		errorMsg := fmt.Sprintf("unknown template '%s'", template)
		logger.PrintError(errorMsg)
	}

	templateFunc()
}
