package main

import (
	"image/png"
	"os"
	"time"

	"github.com/cufee/aftermath-core/internal/core/localization"
	"github.com/cufee/aftermath-core/internal/logic/render"
	"github.com/cufee/aftermath-core/internal/logic/stats"
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

	now := time.Now()
	img, err := render.RenderStatsImage(session, nil, localization.LanguageEN)
	if err != nil {
		panic(err)
	}

	println(time.Since(now).String())

	// now := time.Now()
	// img, err := render.RenderCards(cards, &render.RenderOptions{
	// 	Style: render.Style{
	// 		Font:      render.FontLarge,
	// 		FontColor: render.FontMediumColor,

	// 		Direction:  render.DirectionVertical,
	// 		AlignItems: render.AlignItemsCenter,
	// 		PaddingX:   20,
	// 		PaddingY:   20,
	// 		Gap:        20,
	// 	},
	// 	Debug: false,
	// })
	// if err != nil {
	// 	panic(err)
	// }
	// println(time.Since(now).String())

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
