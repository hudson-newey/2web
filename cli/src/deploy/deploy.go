package deploy

import "github.com/hudson-newey/2web-cli/src/deploy/netlify"

func DeploySolution() {
	location := determineDeploymentLocation()

	switch location {
	case NetlifyLocation:
		netlify.Deploy()
	}
}

func determineDeploymentLocation() deploymentLocation {
	// We currently only support deployments to Netlify
	// TODO: add support for other locations
	return NetlifyLocation
}
