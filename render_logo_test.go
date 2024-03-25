package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"testing"

	"github.com/cufee/aftermath-core/internal/logic/render"
	"github.com/cufee/aftermath-core/internal/logic/render/shared"
)

func TestFullLogoRenderPipeline(t *testing.T) {
	sizes := []float64{32, 64, 128, 256, 512}
	colors := map[string]color.Color{"blue": render.TextSubscriptionPlus, "red": render.ColorAftermathRed, "yellow": render.TextSubscriptionPremium}

	result := make(map[string]image.Image)
	for colorName, color := range colors {
		logo := shared.AftermathLogo(color, shared.LogoSizingOptions{
			Gap:       40,
			Jump:      60,
			Lines:     9,
			LineStep:  120,
			LineWidth: 60,
		})
		for _, size := range sizes {
			{
				withPadding := render.NewBlocksContent(render.Style{
					// PaddingX: size / 10,
					// PaddingY: size / 10,
					// BackgroundColor: render.DiscordBackgroundColor,
				}, render.NewImageContent(render.Style{Width: size, Height: size}, logo))
				logoImage, err := withPadding.Render()
				if err != nil {
					t.Fatal(err)
				}

				// bgImage, _ := assets.GetImage("images/backgrounds/default")
				// img := core.AddBackground(logoImage, bgImage, core.Style{})
				fileName := fmt.Sprintf("logo-%s-%d.png", colorName, int(size))
				result[fileName] = logoImage

				f, err := os.Create("static/" + fileName)
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

			{
				withBackground := render.NewBlocksContent(render.Style{
					PaddingX:        size / 5,
					PaddingY:        size / 5,
					BackgroundColor: render.DiscordBackgroundColor,
				}, render.NewImageContent(render.Style{Width: size, Height: size}, logo))
				withBackgroundImage, err := withBackground.Render()
				if err != nil {
					t.Fatal(err)
				}

				fileName := fmt.Sprintf("logo-%s-%d-dark.png", colorName, int(size))
				result[fileName] = withBackgroundImage

				f, err := os.Create("static/" + fileName)
				if err != nil {
					t.Fatal(err)
				}
				defer f.Close()

				// err = png.Encode(f, img)
				err = png.Encode(f, withBackgroundImage)
				if err != nil {
					t.Fatal(err)
				}
			}
		}
	}

	// var ogFontSize float64 = 96
	// fonts, _ := assets.GetFontFaces("default", ogFontSize)
	// ogFont := fonts[ogFontSize]

	ogBlock := render.NewBlocksContent(render.Style{
		Gap:             30,
		Width:           1000,
		Height:          1000,
		AlignItems:      render.AlignItemsCenter,
		Direction:       render.DirectionVertical,
		JustifyContent:  render.JustifyContentCenter,
		BackgroundColor: render.DiscordBackgroundColor,
	},
		render.NewImageContent(render.Style{}, result["logo-red-512.png"]),
		// render.NewBlocksContent(render.Style{Direction: render.DirectionVertical, AlignItems: render.AlignItemsCenter, JustifyContent: render.JustifyContentCenter},
		// 	render.NewTextContent(render.Style{Font: &ogFont, FontColor: render.TextPrimary, PaddingY: -20}, "Aftermath"),
		// 	// render.NewTextContent(render.Style{Font: &render.Font2XL, FontColor: render.TextPrimary, PaddingY: -10}, "Your Blitz stats, fast and beautiful"),
		// ),
	)

	ogImage, _ := ogBlock.Render()

	{
		f, err := os.Create("static/og.png")
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		err = png.Encode(f, ogImage)
		if err != nil {
			t.Fatal(err)
		}
	}
	{
		f, err := os.Create("static/og.jpg")
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		err = jpeg.Encode(f, ogImage, &jpeg.Options{Quality: 100})
		if err != nil {
			t.Fatal(err)
		}
	}
}
