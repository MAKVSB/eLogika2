package common

type ErrorObject map[string]any

type ErrorResources struct {
	ResourceType string `json:"resourceType"`
	ResourceID   uint   `json:"resourceId"`
	ResourceData string `json:"resourceData,omitempty"`
}

type ErrorFile struct {
	Content  string `json:"content"`
	MimeType string `json:"mimeType"`
}

// @Description Generic error message
type ErrorResponse struct {
	Code       int              `json:"code,omitempty"`
	Message    string           `json:"message,omitempty"`
	Details    any              `json:"details,omitempty"`
	Resources  []ErrorResources `json:"resources,omitempty"`
	FormErrors ErrorObject      `json:"formErrors,omitempty"`
	FileData   *ErrorFile       `json:"fileData,omitempty"`
}
