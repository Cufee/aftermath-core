package session

import (
	"image"

	"github.com/cufee/aftermath-core/internal/logic/render"
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

	return cardsImage, nil
}
