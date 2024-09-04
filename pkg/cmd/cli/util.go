package cli

import (
	"encoding/json"
	"fmt"
)

func dataToJson(data any) string {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println("Error serializing to JSON:", err)
	}
	return string(jsonData)
}
