package pkg

import (
	"encoding/json"
)

type ValidationError map[string]string

func (v ValidationError) Error() string {
	errJson, _ := json.Marshal(v)
	return string(errJson)
}
