package utils

import (
	"bytes"
	"encoding/gob"
)

/*
DeepCopy copies src to dist. dist must be a pointer to a struct.
*/
func DeepCopy[T any](src interface{}, dist *T) (err error) {
	buf := bytes.Buffer{}
	if err = gob.NewEncoder(&buf).Encode(src); err != nil {
		return
	}
	return gob.NewDecoder(&buf).Decode(dist)
}
