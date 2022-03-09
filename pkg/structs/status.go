package structs

import (
	"encoding/json"
	"fmt"
)

type Status struct {
	Name  string
	Type  string
	Value string
	Notes string `json:",omitempty"` // optional, not yet used
}

func (s Status) ToJSON() string {
	b, err := json.Marshal(s)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(b)
}
