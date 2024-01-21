package content

import (
	"bufio"
	"bytes"
	"encoding/base64"

	"fmt"
	"image"
	"image/jpeg"
	"io"
	"net/http"
	"net/url"

	"github.com/cufee/aftermath-core/errors"
	"github.com/cufee/aftermath-core/internal/core/cloudinary"
)

func UploadUserImage(userID, remoteImage string) (string, error) {
	encodedImage, err := EncodeRemoteImage(remoteImage)
	if err != nil {
		return "", err
	}

	link, err := cloudinary.DefaultClient.UploadWithModeration(userID, encodedImage)
	if err != nil {
		return "", err
	}

	return link, nil
}

func EncodeRemoteImage(remoteImage string) (string, error) {
	remoteUrl, err := url.Parse(remoteImage)
	if err != nil {
		return "", err
	}
	remoteUrl.RawQuery = ""

	res, err := http.DefaultClient.Get(remoteUrl.String())
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	rawImage, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	img, format, err := image.Decode(bytes.NewReader(rawImage))
	if err != nil {
		return "", errors.ErrInvalidImageFormat
	}

	var encodedImage string
	switch format {
	case "png":
		fallthrough
	case "jpg":
		fallthrough
	case "jpeg":
		var data bytes.Buffer
		err := jpeg.Encode(bufio.NewWriter(&data), img, &jpeg.Options{Quality: jpeg.DefaultQuality})
		if err != nil {
			return "", err
		}
		encodedImage = fmt.Sprintf("data:image/jpeg;base64,%s", base64.StdEncoding.EncodeToString(data.Bytes()))
	default:
		return "", errors.ErrInvalidImageFormat
	}

	return encodedImage, nil
}
