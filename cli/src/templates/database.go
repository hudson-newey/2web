package templates

import "github.com/hudson-newey/2web-cli/src/sdk"

func DatabaseTemplate() {
	sdk.CopyFromSdk("database", "database")
}
