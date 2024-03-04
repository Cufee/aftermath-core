package main

import (
	"image/png"
	"os"
	"testing"

	"github.com/cufee/aftermath-core/internal/logic/render"
	"github.com/cufee/aftermath-core/internal/logic/render/shared"
)

func TestFullLogoRenderPipeline(t *testing.T) {
	logo := shared.AftermathLogo(render.ColorAftermathBlue, shared.DefaultLogoOptions())
	withPadding := render.NewBlocksContent(render.Style{PaddingX: 20, PaddingY: 20, BackgroundColor: render.DiscordBackgroundColor}, render.NewImageContent(render.Style{Width: 150, Height: 150}, logo))
	logoImage, err := withPadding.Render()
	if err != nil {
		t.Fatal(err)
	}

	// bgImage, _ := assets.GetImage("images/backgrounds/default")
	// img := core.AddBackground(logoImage, bgImage, core.Style{})

	f, err := os.Create("test-logo.png")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	// err = png.Encode(f, img)
	err = png.Encode(f, logoImage)
	if err != nil {
		t.Fatal(err)
	}
}
