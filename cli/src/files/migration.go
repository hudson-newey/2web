package files

import "regexp"

type Migration struct {
	TargetPath  string
	Selector    *regexp.Regexp
	Replacement string
}
