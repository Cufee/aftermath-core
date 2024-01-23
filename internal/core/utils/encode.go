package utils

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/png"
)

func EncodeImage(img image.Image) (string, error) {
	encoded := new(bytes.Buffer)
	err := png.Encode(encoded, img)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(encoded.Bytes()), nil
}
