package content

import "testing"

func TestPickRandomBackgroundImages(t *testing.T) {
	images, err := PickRandomBackgroundImages(3)
	if err != nil {
		t.Error(err)
	}

	t.Log(images)
}
