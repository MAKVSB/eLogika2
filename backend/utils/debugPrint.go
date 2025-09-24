package utils

import (
	"encoding/json"
	"fmt"
)

func DebugPrintJSON(v any) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Println("error marshaling to JSON:", err)
		return
	}
	fmt.Println(string(b))
}
