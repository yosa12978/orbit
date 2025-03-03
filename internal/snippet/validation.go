package snippet

import (
	"regexp"
)

var contentRE = regexp.MustCompile(`^.+$`)

func validateContent(content string) error {
	if !contentRE.MatchString(content) {
		return ErrInvalidContent
	}
	return nil
}
