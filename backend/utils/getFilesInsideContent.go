package utils

import "encoding/json"

func GetFilesInsideContent(jsonContent json.RawMessage) []int {
	var imageIds []int

	var inner func(content json.RawMessage)
	inner = func(content json.RawMessage) {
		var node map[string]json.RawMessage
		if err := json.Unmarshal(content, &node); err != nil {
			return
		}

		// Check if type == "custom-image"
		var nodeType string
		if err := json.Unmarshal(node["type"], &nodeType); err == nil && nodeType == "custom-image" {
			// Extract attrs
			var attrs struct {
				Mode string `json:"mode"`
				ID   int    `json:"id"`
			}
			if err := json.Unmarshal(node["attrs"], &attrs); err == nil {
				if attrs.Mode == "storage" {
					imageIds = append(imageIds, attrs.ID)
				}
			}
		}

		// Recurse into content array if present
		var children []json.RawMessage
		if err := json.Unmarshal(node["content"], &children); err == nil {
			for _, child := range children {
				inner(child)
			}
		}
	}

	inner(jsonContent)
	return imageIds
}
