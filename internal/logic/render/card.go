package render

import (
	"image"
	"image/color"

	core "github.com/cufee/aftermath-core/internal/core/stats"
	"github.com/cufee/aftermath-core/internal/logic/stats"
)

type card struct {
	title  string
	blocks []block

	options     RenderOptions
	titleConfig RenderConfig

	statsOptions RenderOptions

	renderedImages             []image.Image
	renderedContentMaxWidth    int
	renderedContentTotalWidth  int
	renderedContentMaxHeight   int
	renderedContentTotalHeight int
}

var largeStatsCardOptions card = card{
	statsOptions: RenderOptions{
		Style: Style{
			JustifyContent: JustifyContentSpaceEvenly,
			Direction:      DirectionHorizontal,
			AlignItems:     AlignItemsCenter,
			Gap:            5,
			PaddingX:       20,
			PaddingY:       20,

			// Width: 600,
		},
		Debug: false,
	},
	options: RenderOptions{
		Style: Style{
			Direction:  DirectionVertical,
			AlignItems: AlignItemsCenter,

			BackgroundColor: color.RGBA{R: 50, G: 50, B: 50, A: 200},
			BorderRadius:    10,
		},
	},
	titleConfig: defaultBlockRenderConfig.Career,
}

func (card *card) RenderBlocks() ([]image.Image, error) {
	if card.renderedImages != nil {
		return card.renderedImages, nil
	}

	for _, block := range card.blocks {
		img, err := block.Render()
		if err != nil {
			return nil, err
		}

		switch block.options.Style.Direction {
		case DirectionHorizontal:
			card.renderedContentTotalWidth += img.Bounds().Dx()
			card.renderedContentMaxWidth = card.renderedContentTotalWidth

			if img.Bounds().Dy() > card.renderedContentMaxHeight {
				card.renderedContentMaxHeight = img.Bounds().Dy()
				card.renderedContentTotalHeight = card.renderedContentMaxHeight
			}
		case DirectionVertical:
			card.renderedContentTotalHeight += img.Bounds().Dy()
			card.renderedContentMaxHeight = card.renderedContentTotalHeight

			if img.Bounds().Dx() > card.renderedContentMaxWidth {
				card.renderedContentMaxWidth = img.Bounds().Dx()
				card.renderedContentTotalWidth = card.renderedContentMaxWidth
			}
		}
		card.renderedImages = append(card.renderedImages, img)
	}

	return card.renderedImages, nil
}

func (card *card) Render() (image.Image, error) {
	images, err := card.RenderBlocks()
	if err != nil {
		return nil, err
	}

	statsImage, err := renderImages(images, card.statsOptions)
	if err != nil {
		return nil, err
	}

	titleRow := blockRow{
		value:  card.title,
		config: card.titleConfig,
	}

	titleImage, err := titleRow.Render()
	if err != nil {
		return nil, err
	}

	return renderImages([]image.Image{titleImage, statsImage}, card.options)
}

func SnapshotToCards(snapshot *stats.Snapshot, averages *core.ReducedStatsFrame) ([]card, error) {
	var cards []card

	// Player Title

	// Regular Battles
	overviewUnratedCard := largeStatsCardOptions
	overviewUnratedCard.title = "Regular Battles" // TODO: localize
	overviewUnratedCard.blocks = FrameToLargeStatsBlocks(snapshot.Diff.Global, snapshot.Selected.Global, averages, &defaultBlockRenderConfig)
	cards = append(cards, overviewUnratedCard)

	// Rating Battles

	// Vehicles
	for _, vehicle := range snapshot.Diff.Vehicles {
		vehicleCard := largeStatsCardOptions
		vehicleCard.title = "Vehicle" // TODO: localize
		vehicleCard.blocks = FrameToLargeStatsBlocks(vehicle.ReducedStatsFrame, snapshot.Selected.Vehicles[vehicle.VehicleID].ReducedStatsFrame, averages, &defaultBlockRenderConfig)

		cards = append(cards, vehicleCard)
	}

	return cards, nil
}

func RenderCards(cards []card, options RenderOptions) (image.Image, error) {
	maxContentWidth, maxContentHeight := 0, 0
	// totalContentWidth, totalContentHeight := 0, 0

	// totalPaddingX, totalPaddingY, totalGap := 0.0, 0.0, 0.0

	var images []image.Image
	for _, card := range cards {
		_, err := card.RenderBlocks()
		if err != nil {
			return nil, err
		}

		if card.renderedContentMaxHeight > maxContentHeight {
			maxContentHeight = card.renderedContentMaxHeight
		}
		if card.renderedContentMaxWidth > maxContentWidth {
			maxContentWidth = card.renderedContentMaxWidth
		}
	}

	for _, card := range cards {
		switch options.Style.Direction {
		case DirectionVertical:
			// card.options.Style.Width = float64(maxContentWidth * len(card.blocks))
		default: // DirectionHorizontal
			// card.options.Style.Height = float64(maxContentHeight * len(card.blocks))
		}

		img, err := card.Render()
		if err != nil {
			return nil, err
		}
		images = append(images, img)
	}

	return renderImages(images, options)
}
