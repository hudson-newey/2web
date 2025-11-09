package reactiveCompiler

import (
	"fmt"
	"hudson-newey/2web/src/content/document"
	"hudson-newey/2web/src/content/javascript"
	"hudson-newey/2web/src/content/page"
	"hudson-newey/2web/src/models"
	"strings"
)

// Warning: static prop variables E.g. *value="$count" is typically an escape
// hatch out of 2web reactivity and optimization.
// It will literally create code to set the underlying property using JavaScript
// as soon as the website is loaded.
// You should always ask if it is possible to complete your task without the
// star prefix.
// E.g. value="$count" will ONLY be evaluated at compile time and results in an
// SSG site without any runtime overhead.
func compileStaticPropVar(
	pageModel *page.Page,
	varNode *models.ReactiveVariable,
) {
	uniquePropSelectors := getUniqueSelectors(varNode.Props)

	content := pageModel.Html.Content

	for _, propNode := range uniquePropSelectors {
		// Some of the properties can be evaluated at compile time to prevent
		// creating runtime JavaScript.
		// This is a part of my goal to aggressively optimize for SSG.
		if propNode.PropName == "innerText" {
			returnedContent, err := replaceStaticTextProp(content, varNode, propNode)
			if err == nil {
				content = returnedContent
				continue
			}
		}

		elementSelector := javascript.CreateJsElementName()

		// Preserve the original node selector so that other reactivity classes can
		// target this element.
		content = strings.ReplaceAll(content, propNode.Node.Selector, propNode.Node.Selector+" "+elementSelector)

		// If there are multiple instances of the same reactive property
		// selector, we can combine the assignment into one querySelectorAll
		// this means that we don't have to do multiple properties to update
		// elements that have the same selector.
		htmlSource := ""
		selectorCount := strings.Count(content, propNode.Node.Selector)
		if selectorCount > 1 {
			// We use forEach here instead of map() so that the JavaScript
			// JIT doesn't have to keep track of a return value.
			//
			// We use the square brackets here because some properties have dashes which
			// cannot be acceded with a period.
			//
			// There's a newline at the start of this script tag so that when it is
			// appended to the body, it's on its own line, and semantically distinct.
			htmlSource = fmt.Sprintf(`
				<script type="module">
					document.querySelectorAll("[%s]").forEach((__2_element_ref_mod) => __2_element_ref_mod["%s"] = %s);
				</script>`,
				elementSelector, propNode.PropName, propNode.Variable.InitialValue,
			)
		} else {
			htmlSource = fmt.Sprintf(`
				<script type="module">
					document.querySelector("[%s]")["%s"] = %s;
				</script>`,
				elementSelector, propNode.PropName, propNode.Variable.InitialValue,
			)
		}

		injectableTemplate, err := document.BuildTemplate(htmlSource, nil)
		if err != nil {
			panic(err)
		}

		// If there is no body tag, we just append to the end of the content.
		// TODO: This should be improved once we start using lazy loaded scripts
		// instead of inlining reactivity as scripts in HTML content.
		if document.HasBodyTag(content) {
			content = document.InjectContent(content, injectableTemplate, document.BodyTop)
		} else {
			content = document.InjectContent(content, injectableTemplate, document.Leading)
		}
	}

	pageModel.Html.Content = content
}

// If a node has an *innerText attribute and the underlying variable is static,
// then the *innerText can be evaluated at compile time instead of creating
// runtime javascript.
//
// E.g.
//
// <script compiled>
// $ message = "Hello World!";
// </script>
// <h1 *innerText="$message"></h1>
//
// In this example, the $message variable is never changed, and we can therefore
// evaluate the innerText at compile time.
// This allows transparent SSG without the programmer having to specify what
// type of rendering they want for the "$message" variable.
func replaceStaticTextProp(
	content string,
	varNode *models.ReactiveVariable,
	propNode *models.ReactiveProperty,
) (string, error) {
	selectorIndex := strings.Index(content, propNode.Node.Selector)
	if selectorIndex == -1 {
		return content, fmt.Errorf("unknown")
	}

	maxLen := len(content)
	openingBracketIndex := selectorIndex

	for openingBracketIndex < maxLen {
		searchChar := content[openingBracketIndex]
		openingBracketIndex++

		if searchChar == '>' {
			break
		} else if searchChar == '/' {
			// We cannot evaluate the innerText at runtime because the element is a
			// self-closing tag.
			// Therefore, we fallback to setting the underlying innerText property
			// and let the custom-element / DOM handle the behavior
			return content, fmt.Errorf("attempted to inline self-closing tag")
		}
	}

	closingBracketIndex := openingBracketIndex
	for closingBracketIndex < maxLen-1 {
		searchString := content[closingBracketIndex : closingBracketIndex+2]
		if searchString == "</" {
			break
		} else if searchString[0] == '<' {
			// We have multiple nested children under this element, and I don't want
			// to code out the cases to handle this condition. Fallback to runtime
			// assignment
			// TODO: handle this edge case
			return content, fmt.Errorf("innerText element has children")
		}

		closingBracketIndex++
	}

	initialTextValue := strings.TrimPrefix(varNode.InitialValue, "\"")
	initialTextValue = strings.TrimSuffix(initialTextValue, "\"")
	initialTextValue = strings.TrimPrefix(initialTextValue, "'")
	initialTextValue = strings.TrimSuffix(initialTextValue, "'")

	// Replace the inner text of the element with the variables initial value
	// we do not have to worry about cleaning up the attribute because it will
	// be handled later final step of the compiler.
	content = replaceBetween(content, openingBracketIndex, closingBracketIndex, initialTextValue)

	return content, nil
}

func replaceBetween(s string, start int, end int, replacement string) string {
	return s[:start] + replacement + s[end:]
}
