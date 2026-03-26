package input

import (
	"errors"
	"strings"
)

func filterKeys(raw string, key_map map[string]string) error {

	// This avoids string being empty, because it is mandatory to have ":"
	if !strings.Contains(raw, ":") {
		return errors.New("Invalid key! -> " + raw)
	}
	key := strings.SplitN(raw, ":", 2)
	if key[0] == "" || key[1] == "" {
		return errors.New("Invalid key! -> " + raw)
	}

	key[0] = strings.TrimSpace(key[0])
	key[1] = strings.TrimSpace(key[1])
	key_map[key[0]] = key[1] // key_map[key_name] = key_value
	return nil
}
