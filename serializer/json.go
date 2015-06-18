package serializer

import (
	"encoding/json"
)

type Json struct {
}

func (s Json) Marshal(resp []interface{}) (data []byte, err error) {
	data, err = json.Marshal(resp)
	return
}

func (s Json) Unmarshal(data []byte) (msg []interface{}, err error) {
	err = json.Unmarshal(data, &msg)
	// TODO: check msg format
	return
}
