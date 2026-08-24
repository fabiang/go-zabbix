package types

import (
	"encoding/json"
	"fmt"
)

type ZBXBoolean bool

func (bit *ZBXBoolean) UnmarshalJSON(data []byte) error {
	var str string
	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}

	switch str {
	case "1", "true":
		*bit = true

	case "0", "false", "":
		*bit = false

	default:
		return fmt.Errorf("Boolean unmarshal error. Invalid input %q", str)
	}

	return nil
}

func (bit ZBXBoolean) MarshalJSON() ([]byte, error) {
	if bit == true {
		return json.Marshal("1")
	}

	return json.Marshal("0")
}
