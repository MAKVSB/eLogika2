package models

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"gorm.io/gorm/schema"
)

type Mark struct {
	Type  string                 `json:"type"`
	Attrs map[string]interface{} `json:"attrs,omitempty"`
}

type TipTapContent struct {
	Type    string                 `json:"type,omitempty"`
	Attrs   map[string]interface{} `json:"attrs,omitempty"`
	Content []*TipTapContent       `json:"content,omitempty"`
	Marks   []Mark                 `json:"marks,omitempty"`
	Text    string                 `json:"text,omitempty"`
}

func (ttc TipTapContent) Hash() (string, error) {
	data, err := json.Marshal(ttc)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:]), nil
}

func (ttc *TipTapContent) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	result := TipTapContent{}
	err := json.Unmarshal(bytes, &result)
	*ttc = TipTapContent(result)
	return err
}

func (ttc *TipTapContent) Value(ctx context.Context, field *schema.Field, dst reflect.Value, fieldValue interface{}) (interface{}, error) {
	if ttc == nil {
		return []byte{}, nil
	}
	return json.Marshal(ttc)
}
