package main

import (
	"image/png"
	"os"
	"testing"

	"github.com/cufee/aftermath-core/internal/core/database"
	"github.com/cufee/aftermath-core/internal/core/utils"
	"github.com/cufee/aftermath-core/internal/logic/preview"
)

func TestPreviewRenderPipeline(t *testing.T) {
	err := database.Connect(utils.MustGetEnv("DATABASE_URL"))
	if err != nil {
		t.Fatal(err)
	}

	data, err := database.GetAppConfiguration[[]string]("backgroundImagesSelection")
	if err != nil {
		t.Fatal(err)
	}
	image, err := preview.RenderBackgroundPreview("Your Awesome Nickname", "", data.Value)
	if err != nil {
		t.Fatal(err)
	}

	f, err := os.Create("test-preview.png")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	err = png.Encode(f, image)
	if err != nil {
		t.Fatal(err)
	}
}
