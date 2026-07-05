package debugjson

import (
	"fmt"
	"hudson-newey/2web/src/cli"
	jsonFile "hudson-newey/2web/src/content/json"
	"hudson-newey/2web/src/content/page"
	"hudson-newey/2web/src/debugger"
	"strings"
)

func GenerateDebugJson() {
	file := jsonFile.JsonFile{}
	file.AddContent(debuggerInfoString())

	pageModel := page.Page{
		Html: &file,
	}

	outPath := cli.GetEnvVars().DebugOverride
	pageModel.WriteHtml(outPath)
}

func debuggerInfoString() string {
	allDebugInfo := debugger.GetDebugInfo()

	var messagesSb strings.Builder
	for i, info := range allDebugInfo {
		entry := fmt.Sprintf("'%s'", info.Message)
		messagesSb.WriteString(entry)
		if i != 0 {
			messagesSb.WriteString(",")
		}
	}

	return fmt.Sprintf(`{
  "messages": [%s]
}`, messagesSb.String())
}
