package main

import (
	"image/color"
	"image/png"
	"os"

	"github.com/cufee/aftermath-core/internal/core/localization"
	"github.com/cufee/aftermath-core/internal/logic/render"
	"github.com/cufee/aftermath-core/internal/logic/render/assets"
	"github.com/cufee/aftermath-core/internal/logic/stats"
)

var (

	// debugColorRed   = color.RGBA{255, 0, 0, 255}
	// debugColorPink = color.RGBA{255, 192, 203, 255}
	// debugColorGreen = color.RGBA{20, 160, 20, 255}

	// FontSizeHeader   = 36.0
	// FontSize         = 24.0
	// TextMargin       = FontSize / 2
	// FrameWidth       = 900
	// FrameMargin      = 50
	// BaseCardWidth    = FrameWidth - (2 * FrameMargin)
	// BaseCardHeigh    = 150
	// BaseCardColor    = color.RGBA{30, 30, 30, 204}
	// DecorLinesColor  = color.RGBA{80, 80, 80, 255}

	BigTextColor   = color.RGBA{255, 255, 255, 255} // Session stats values, titles and names
	SmallTextColor = color.RGBA{204, 204, 204, 255} // Career stats values
	AltTextColor   = color.RGBA{100, 100, 100, 255} // Stats labels
	// PremiumColor  = color.RGBA{255, 223, 0, 255} // Premium Vehicle
	// VerifiedColor = color.RGBA{72, 167, 250, 255} // Verified Account
)

func main() {
	// err := cache.RefreshSessions(cache.SessionTypeDaily, "na", 1013072123)
	// if err != nil {
	// 	panic(err)
	// }

	session, err := stats.GetCurrentPlayerSession("na", 1013072123) // 1013379500 1013072123 1039212441
	if err != nil {
		panic(err)
	}

	fontFaces, ok := assets.GetFontFaces("default", 24, 18, 14)
	if !ok {
		panic("font not found")
	}

	config := render.BlockRenderConfig{
		Session: render.RenderConfig{
			Font:      fontFaces[24],
			FontColor: BigTextColor,
			// Debug:     true,
		},
		Career: render.RenderConfig{
			Font:      fontFaces[18],
			FontColor: SmallTextColor,
			// Debug:     true,
		},
		Label: render.RenderConfig{
			Font:      fontFaces[14],
			FontColor: AltTextColor,
			// Debug:     true,
		},
		RowStyle: render.Style{
			Direction:  render.DirectionVertical,
			AlignItems: render.AlignItemsCenter,
			Gap:        5,
		},
		SetStyle: render.Style{
			Direction:  render.DirectionHorizontal,
			AlignItems: render.AlignItemsCenter,
			Gap:        20,
		},
		Locale: localization.LanguageEN,
	}

	blocks := render.FrameToLargeStatsBlocks(session.Diff.Global, session.Selected.Global, nil, &config)
	img, err := blocks.Render()
	if err != nil {
		panic(err)
	}

	f, err := os.Create("test.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	err = png.Encode(f, img)
	if err != nil {
		panic(err)
	}
}
