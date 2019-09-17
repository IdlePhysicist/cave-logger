package common

import (
	"bytes"
	"encoding/json"
)

// StructToJSON convert struct to json.
func StructToJSON(i interface{}) string {
	j, err := json.Marshal(i)
	if err != nil {
		return ""
	}

	out := new(bytes.Buffer)
	json.Indent(out, j, "", "    ")
	return out.String()
}