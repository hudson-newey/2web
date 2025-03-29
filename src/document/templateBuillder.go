package document

import (
	"bytes"
	"html/template"
)

func BuildTemplate(source string, data any) (string, error) {
	errorTemplate := template.Must(template.New("errorTemplate").Parse(source))

	var buf bytes.Buffer
	err := errorTemplate.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
