package documentErrors

import (
	"bytes"
	"html/template"
	"hudson-newey/2web/src/models"
)

var errorList []models.Error

func AddError(errorModel models.Error) {
	errorList = append(errorList, errorModel)
}

func InjectErrors(pageContent string) string {
	if len(errorList) == 0 {
		return pageContent
	}

	return pageContent + createErrorTemplate(errorList)
}

// creates a HTML error template that can be used to display errors
// in the browser
func createErrorTemplate(errors []models.Error) string {
	htmlTemplate := `
        <div>
            <div class="__2_error_container">
                <h1 class="__2_error_header">
                    <strong>Compiler Error</strong>
                </h1>
                <div class="__2_error_list">
                    {{range .}}
                        <div class="__2_error">
                            <h2 class="__2_error_file">{{.FilePath}}</h2>
                            <pre class="__2_error_message">{{.Message}}</pre>
                        </div>
                    {{end}}
                </div>
            </div>

            <style>
                /* make the rest of the page a dark opaque color */
                html {
                    background-color: rgba(0, 0, 0, 0.5);
                }

                .__2_error_container {
                    position: fixed;
                    inset: 5%;
                    padding: 2rem 5rem;
                    z-index: 5000;

                    background: radial-gradient(circle, rgba(5, 0, 0, 0.9) 0%, rgba(20, 0, 0, 0.98) 99%, rgba(40, 0, 0, 1) 100%);
                    background-color: rgba(20, 0, 0);
                    border-radius: 1rem;

                    box-shadow: 0 0.5rem 10rem rgba(0, 0, 0);
                }

                .__2_error_header {
                    font-weight: bold;
                    font-size: 2.2rem;

                    color: #fff;
                    margin-bottom: 2rem;
                    border-bottom: 1px solid #fff;
                    padding-bottom: 1rem;
                }

                .__2_error_list {
                    & > .__2_error {
                        color: white;
                        margin-bottom: 1.5rem;
                        padding: 1rem;
                        background-color: rgba(80, 20, 20, 0.1);
                        border-radius: 0.5rem;
                        box-shadow: 0 0 0.25rem rgba(149, 38, 43, 0.9);
                    }

                    .__2_error_file {
                        color: #f00;
                        font-weight: bold;
                        margin-bottom: 1rem;
                        margin-top: 0;
                    }

                    .__2_error_message {
                        color: rgb(255, 230, 230);
                        margin: 0;
                        font-size: 1.1rem;
                        font-weight: 400;
                        white-space: pre-wrap;
                        line-height: 2;
                        
                        overflow-wrap: break-word;
                    }
                }
            </style>
        <div>
    `

	errorTemplate := template.Must(template.New("errorTemplate").Parse(htmlTemplate))
	var buf bytes.Buffer
	err := errorTemplate.Execute(&buf, errors)
	if err != nil {
		// Handle the error, maybe add it to errorList
		AddError(models.Error{Message: "Failed to render error template: " + err.Error()})
	}
	errorHtml := buf.String()

	return errorHtml
}
