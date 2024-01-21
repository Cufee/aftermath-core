package content

import (
	"math/rand"

	"github.com/cufee/aftermath-core/internal/core/cloudinary"
)

func PickRandomBackgroundImages(number int) ([]string, error) {
	if number <= 0 {
		return nil, nil
	}

	images, err := cloudinary.DefaultClient.GetFolderImages("Aftermath/manual-uploads", 10)
	if err != nil {
		return nil, err
	}
	if len(images) < number || len(images) == 0 {
		return images, nil
	}

	rand.Shuffle(len(images), func(i, j int) { images[i], images[j] = images[j], images[i] })
	return images[:number], nil
}
