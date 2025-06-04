package templates

import "github.com/hudson-newey/2web-cli/src/sdk"

func LoadBalancerTemplate() {
	sdk.CopyFromSdk("load-balancer", "load-balancer")
}
