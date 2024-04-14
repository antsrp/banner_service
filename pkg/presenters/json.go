package mapper

import (
	"encoding/json"
	"fmt"
)

type Indent struct {
	prefix string
	indent string
}

var DefaultIndent = Indent{prefix: "", indent: "\t"}

func NewIndent(prefix, indent string) *Indent {
	return &Indent{
		prefix: prefix,
		indent: indent,
	}
}

func ToJSON[T any](data T, settings *Indent) ([]byte, error) {
	result, err := json.MarshalIndent(data, settings.prefix, settings.indent)
	if err != nil {
		return nil, fmt.Errorf("can't marshal json: %w", err)
	}
	return result, nil
}

func FromJSON[T any](data []byte) (*T, error) {
	var result T
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("can't unmarshal json: %w", err)
	}
	return &result, nil
}
