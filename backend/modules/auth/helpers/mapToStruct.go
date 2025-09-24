package helpers

import (
	"encoding/json"

	"elogika.vsb.cz/backend/modules/auth/dtos"
)

func MapToStruct(input map[string]interface{}) (*dtos.LoggedUserDTO, error) {
	// Marshal to JSON
	bytes, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	// Unmarshal to struct
	var user dtos.LoggedUserDTO
	err = json.Unmarshal(bytes, &user)
	return &user, err
}
