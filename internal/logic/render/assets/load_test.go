package assets

import (
	"testing"
)

func TestLoadAllFiles(t *testing.T) {
	files, err := getAllFiles(imagesEmbed, ".")
	if err != nil {
		t.Fatal(err)
	}
	if len(files) == 0 {
		t.Fatal("no files loaded")
	}
}
