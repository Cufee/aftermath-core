package session

import (
	"image"
	"image/color"

	"github.com/cufee/aftermath-core/internal/logic/render"
	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
)

func RenderStatsImage(player PlayerData, options RenderOptions) (image.Image, error) {
	cards, err := snapshotToCardsBlocks(player, options)
	if err != nil {
		return nil, err
	}

	allCards := render.NewBlocksContent(
		render.Style{
			Direction:  render.DirectionVertical,
			AlignItems: render.AlignItemsCenter,
			PaddingX:   20,
			PaddingY:   20,
			Gap:        10,
			// Debug:      true,
		}, cards...)

	cardsImage, err := allCards.Render()
	if err != nil {
		return nil, err
	}

	// Fill background with black and round the corners
	frameCtx := gg.NewContextForImage(cardsImage)
	frameCtx.SetColor(color.RGBA{44, 47, 51, 255})
	frameCtx.Clear()
	frameCtx.DrawRoundedRectangle(0, 0, float64(frameCtx.Width()), float64(frameCtx.Height()), 20)
	frameCtx.Clip()

	// Resize the background image to fit the cards
	bgImage := imaging.Fill(options.BackgroundImage, frameCtx.Width(), frameCtx.Height(), imaging.Center, imaging.NearestNeighbor)
	bgImage = imaging.Blur(bgImage, 10)
	frameCtx.DrawImage(bgImage, 0, 0)
	frameCtx.DrawImage(cardsImage, 0, 0)

	return frameCtx.Image(), nil
}
