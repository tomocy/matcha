package infra

import (
	"encoding/json"
	"io"
)

func readJSON(src io.Reader, dest interface{}) error {
	return json.NewDecoder(src).Decode(dest)
}
