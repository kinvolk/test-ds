package main

import (
	"encoding/json"
	"fmt"
)

func Dump(thing interface{}) string {
	encoded, err := json.MarshalIndent(thing, "", "  ")
	if err != nil {
		return fmt.Sprintf("<FAILED TO MARSHAL %T: %v>", thing, err)
	}
	return (string)(encoded)
}
