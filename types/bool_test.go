package types

import (
	"encoding/json"
	"testing"
)

func TestBool(t *testing.T) {
	tests := map[string]bool{
		"1":     true,
		"true":  true,
		"0":     false,
		"false": false,
	}

	for input, expected := range tests {
		var current ZBXBoolean
		jsonInput, _ := json.Marshal(input)
		err := json.Unmarshal(jsonInput, &current)
		if err != nil {
			t.Error(err)
		}

		if current != ZBXBoolean(expected) {
			t.Errorf("Expected %q to be %t", input, expected)
		}
	}
}
