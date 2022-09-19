package structs

import (
	"encoding/json"
	"fmt"
)

type Suite struct {
	Name        string
	Description string
	Owner       string
	Notes       string
}


func (s Suite) ToJSON() string {
	b, err := json.Marshal(s)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(b)
}
