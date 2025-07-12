package cfFor

import "strings"

// you can use the template token to use the array value in your template
// e.g.
// {% for 1,2,3 <li>{{$value}}</li> %}
const listSeparatorToken = ","
const replacementToken = "{{&value}}"

func ForLoopContent(rawValues string, contentTokens []string) string {
	values := strings.Split(rawValues, listSeparatorToken)

	joinedContent := strings.Join(contentTokens, " ")

	templateResult := ""
	for _, value := range values {
		templateResult += strings.ReplaceAll(joinedContent, replacementToken, value)
	}

	return templateResult
}
