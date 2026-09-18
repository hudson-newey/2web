package templates

func TailwindTemplate() {
	copyFromTemplates("tailwind.config.js", "tailwind.config.js")
	copyFromTemplates("postcss.config.js", "postcss.config.js")
}
