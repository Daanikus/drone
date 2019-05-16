package util

import (
	"bytes"
	"encoding/json"
)

func DumpObject(obj interface{}) string {
	buff := &bytes.Buffer{}
	enc := json.NewEncoder(buff)
	enc.Encode(obj)
	return buff.String()
}
