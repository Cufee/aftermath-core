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
	"github.com/cufee/aftermath-core/internal/core/database"
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

func GetCurrentImageSelection() ([]string, error) {
	config, err := database.GetAppConfiguration[[]string]("backgroundImagesSelection")
	if err != nil {
		return nil, err
	}
	return config.Value, nil
}

func LoadRemoteImage(remoteImage string) (image.Image, string, error) {
	remoteUrl, err := url.Parse(remoteImage)
	if err != nil {
		return nil, "", err
	}

	res, err := http.DefaultClient.Get(remoteUrl.String())
	if err != nil {
		return nil, "", err
	}
	defer res.Body.Close()

	rawImage, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, "", err
	}

	img, format, err := image.Decode(bytes.NewReader(rawImage))
	if err != nil {
		return nil, "", err
	}

	return img, format, nil
}

func EncodeRemoteImage(remoteImage string) (string, error) {
	img, format, err := LoadRemoteImage(remoteImage)
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
