package structs

import (
	"encoding/json"
	"fmt"
)

type Stat struct {
	TestRunName string
	Count       int
}

func (s Stat) ToJSON() string {
	b, err := json.Marshal(s)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(b)
}
