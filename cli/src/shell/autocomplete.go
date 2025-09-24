package shell

import (
	"fmt"
	"os"
)

// To use this function, you need to run the following command in your terminal:
// source <(2web --generate-autocomplete-script)
func EchoAutocompleteScript() {
	shellScript := ``

	fmt.Print(shellScript)
}

func InstallAutocomplete() {
	// Write to .bashrc and .zshrc
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting user home directory:", err)
		return
	}

	initScript := "\n# 2web autocomplete\nsource <(2web --generate-autocomplete-script)\n"

	bashrcPath := fmt.Sprintf("%s/.bashrc", homeDir)
	if file, err := os.OpenFile(bashrcPath, os.O_APPEND|os.O_WRONLY, 0644); err == nil {
		defer file.Close()
		if _, err := file.WriteString(initScript); err != nil {
			fmt.Println("Error writing to .bashrc:", err)
			return
		}
	}

	zshrcPath := fmt.Sprintf("%s/.zshrc", homeDir)
	if file, err := os.OpenFile(zshrcPath, os.O_APPEND|os.O_WRONLY, 0644); err == nil {
		defer file.Close()
		if _, err := file.WriteString(initScript); err != nil {
			fmt.Println("Error writing to .zshrc:", err)
			return
		}
	}

	fmt.Println("Autocomplete installed. Please restart your terminal or run 'source ~/.bashrc' or 'source ~/.zshrc'.")
}
