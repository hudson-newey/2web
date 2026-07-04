package cache

import "hudson-newey/2web/src/cli"

func cacheLocation() string {
	return cli.GetEnvVars().CacheOverride
}
