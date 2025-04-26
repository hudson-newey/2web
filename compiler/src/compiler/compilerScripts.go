package compiler

// TODO: the implementation of this function is horrible
func removeCompilerScripts(content string) string {
	startToken := compilerStartToken[0]
	endToken := compilerEndToken[0]

	// If someone creates a really small document, there is no way that it
	// contains any start or end tokens.
	// Therefore, we can short circuit and not deal with a lot of index out of
	// bounds errors.
	if len(content) < len(startToken)+len(endToken) {
		return content
	}

	resultContent := ""
	inCompilerScript := false

	for i := range content {
		if inCompilerScript {
			if i+len(endToken) <= len(content) {
				reachedEndToken := content[i-len(endToken):i] == endToken
				inCompilerScript = !reachedEndToken
			}

		} else {
			if i+len(startToken) <= len(content) {
				reachedStartToken := content[i:i+len(startToken)] == startToken
				inCompilerScript = reachedStartToken
			}
		}

		if !inCompilerScript {
			resultContent += string(content[i])
		}
	}

	return resultContent
}
