package checkbox

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

const enterKey = '\n'
const exitKey = 'q'

type Checkbox struct {
	Message string
	Options Options
}

type Options struct {
	Checked bool
}

func CheckboxGroup(checkboxes []Checkbox) {
	for _, item := range checkboxes {
		singleCheckbox(item.Message, item.Options)
	}
}

func singleCheckbox(message string, options Options) (bool, error) {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return false, err
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	checked := options.Checked

	fmt.Println("(Space to toggle, q to quit)")

	fmt.Printf("\r[ ] %s: ", message)

	buf := make([]byte, 1)
	for {
		_, err := os.Stdin.Read(buf)
		if err != nil {
			return false, err
		}

		switch buf[0] {
		case ' ':
			checked = !checked
			// Move cursor to start of line and print updated checkbox
			fmt.Print("\r")
			if checked {
				fmt.Printf("[x] %s ", message)
			} else {
				fmt.Printf("[ ] %s ", message)
			}
		case enterKey:
			return checked, nil
		case exitKey:
			os.Exit(0)
		}
	}
}
