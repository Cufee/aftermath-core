package period

import (
	"image"

	dataprep "github.com/cufee/aftermath-core/dataprep/period"
	"github.com/cufee/aftermath-core/internal/core/database/models"
	"github.com/cufee/aftermath-core/internal/logic/render"
	"github.com/cufee/aftermath-core/internal/logic/stats/period"
)

type PlayerData struct {
	Stats *period.PeriodStats
	Cards dataprep.Cards

	Subscriptions []models.UserSubscription
}

type RenderOptions struct {
	PromoText []string
	CardStyle render.Style
}

func RenderImage(player PlayerData, options RenderOptions) (image.Image, error) {
	cards, err := generateCards(player, options)
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
