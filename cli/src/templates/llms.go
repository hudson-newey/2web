package templates

func LlmsTemplate() {
	copyFromTemplates("llms.txt", "src/llms.txt")
	copyFromTemplates("llms.txt", "src/llms-full.txt")
}
