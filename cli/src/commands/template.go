package commands

import (
	"fmt"

	"github.com/hudson-newey/2web-cli/src/cli"
	"github.com/hudson-newey/2web-cli/src/templates"
)

func template(programName string, command string, args []string) {
	argsLen := len(args)
	if argsLen < 3 {
		errorMsg := fmt.Sprintf("invalid arguments:\n\texpected: %s %s <template>", programName, command)
		cli.PrintError(errorMsg)
	}

	template := args[2]

	switch template {
	case "serer-side-rendering", "ssr":
		templates.SsrTemplate()
	case "database", "db":
		templates.DatabaseTemplate()
	case "load-balancer", "lb":
		templates.LoadBalancerTemplate()
	case "sitemap", "sitemap.xml":
		templates.SitemapTemplate()
	case "robots.txt":
		templates.RobotsTxtTemplate()
	case "llms.txt":
		templates.LlmsTemplate()
	}
}
