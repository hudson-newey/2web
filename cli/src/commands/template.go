package commands

import (
	"fmt"

	"github.com/hudson-newey/2web-cli/src/cli"
	"github.com/hudson-newey/2web-cli/src/templates"
)

func templateCommand(programName string, command string, args []string) {
	argsLen := len(args)
	if argsLen < 3 {
		errorMsg := fmt.Sprintf("invalid arguments:\n\texpected: %s %s <template>", programName, command)
		cli.PrintError(errorMsg)
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
	}

	templateFunc, exists := templates[template]
	if !exists {
		errorMsg := fmt.Sprintf("unknown template '%s'", template)
		cli.PrintError(errorMsg)
	}

	templateFunc()
}
