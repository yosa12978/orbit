package pkg

import (
	"fmt"
	"strings"
)

type ValidationError map[string]string

func (e ValidationError) Error() string {
	errs := make([]string, 0, len(e))
	for k, v := range e {
		errs = append(errs, fmt.Sprintf("%s - %s", k, v))
	}
	return strings.Join(errs, "\n")
}
