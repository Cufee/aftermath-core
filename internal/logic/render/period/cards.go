package period

import (
	"errors"
	"strings"

	"github.com/cufee/aftermath-core/dataprep"
	"github.com/cufee/aftermath-core/dataprep/period"
	"github.com/cufee/aftermath-core/internal/core/database/models"
	"github.com/cufee/aftermath-core/internal/logic/render"
	"github.com/cufee/aftermath-core/internal/logic/render/badges"
	"github.com/cufee/aftermath-core/internal/logic/render/shared"

	helpers "github.com/cufee/aftermath-core/internal/core/utils"
	"github.com/cufee/aftermath-core/utils"
	"github.com/rs/zerolog/log"
)

func generateCards(player PlayerData, options RenderOptions) ([]render.Block, error) {
	if len(player.Cards.Overview.Blocks) == 0 && len(player.Cards.Highlights) == 0 {
		log.Error().Msg("player cards slice is 0 length, this should not happen")
		return nil, errors.New("no cards provided")
	}

	// Calculate minimal card width to fit all the content
	var cardWidth float64
	overviewColumnWidth := float64(shared.DefaultLogoOptions().Width())
	{
		{
			titleStyle := shared.DefaultPlayerTitleStyle(titleCardStyle(cardWidth))
			clanSize := render.MeasureString(player.Stats.Clan.Tag, *titleStyle.ClanTag.Font)
			nameSize := render.MeasureString(player.Stats.Account.Nickname, *titleStyle.Nickname.Font)
			cardWidth = helpers.Max(cardWidth, titleStyle.TotalPaddingAndGaps()+nameSize.TotalWidth+clanSize.TotalWidth*2)
		}
		{
			rowStyle := getOverviewStyle(cardWidth)
			for _, column := range player.Cards.Overview.Blocks {
				for _, block := range column {
					valueStyle, labelStyle := rowStyle.block(block)

					label := block.Label
					if block.Tag == dataprep.TagWN8 {
						label = shared.GetWN8TierName(int(block.Data.Value))
					}
					labelSize := render.MeasureString(label, *labelStyle.Font)
					valueSize := render.MeasureString(block.Data.String, *valueStyle.Font)

					overviewColumnWidth = helpers.Max(overviewColumnWidth, helpers.Max(labelSize.TotalWidth, valueSize.TotalWidth)+(rowStyle.container.PaddingX*2))
				}
			}

			cardStyle := overviewCardStyle(cardWidth)
			paddingAndGaps := (cardStyle.PaddingX+rowStyle.container.PaddingX+rowStyle.blockContainer.PaddingX)*2 + float64(len(player.Cards.Overview.Blocks)-1)*(cardStyle.Gap+rowStyle.container.Gap+rowStyle.blockContainer.Gap)

			overviewCardContentWidth := overviewColumnWidth * float64(len(player.Cards.Overview.Blocks))
			cardWidth = helpers.Max(cardWidth, overviewCardContentWidth+paddingAndGaps)
		}

		{
			highlightStyle := highlightCardStyle(defaultCardStyle(0))
			var highlightBlocksMaxCount, highlightTitleMaxWidth, highlightBlockMaxSize float64
			for _, highlight := range player.Cards.Highlights {
				// Title and tank name
				metaSize := render.MeasureString(highlight.Meta, *highlightStyle.cardTitle.Font)
				titleSize := render.MeasureString(highlight.Title, *highlightStyle.tankName.Font)
				highlightTitleMaxWidth = helpers.Max(highlightTitleMaxWidth, metaSize.TotalWidth, titleSize.TotalWidth)

				// Blocks
				highlightBlocksMaxCount = helpers.Max(highlightBlocksMaxCount, float64(len(highlight.Blocks)))
				for _, block := range highlight.Blocks {
					labelSize := render.MeasureString(block.Label, *highlightStyle.blockLabel.Font)
					valueSize := render.MeasureString(block.Data.String, *highlightStyle.blockValue.Font)
					highlightBlockMaxSize = helpers.Max(highlightBlockMaxSize, valueSize.TotalWidth, labelSize.TotalWidth)
				}
			}

			highlightCardWidthMax := (highlightStyle.container.PaddingX * 2) + (highlightStyle.container.Gap * highlightBlocksMaxCount) + (highlightBlockMaxSize * highlightBlocksMaxCount) + highlightTitleMaxWidth
			cardWidth = helpers.Max(cardWidth, highlightCardWidthMax)
		}
	}

	var cards []render.Block

	// We first footer in order to calculate the minimum required width
	// Footer Card
	var footerCard render.Block
	{
		var footer []string
		switch strings.ToLower(utils.RealmFromAccountID(player.Stats.Account.ID)) {
		case "na":
			footer = append(footer, "North America")
		case "eu":
			footer = append(footer, "Europe")
		case "as":
			footer = append(footer, "Asia")
		}

		sessionTo := player.Stats.End.Format("January 2, 2006")
		sessionFrom := player.Stats.Start.Format("January 2, 2006")
		if sessionFrom == sessionTo {
			footer = append(footer, sessionTo)
		} else {
			footer = append(footer, sessionFrom+" - "+sessionTo)
		}
		footerBlock := shared.NewFooterCard(strings.Join(footer, " • "))
		footerImage, err := footerBlock.Render()
		if err != nil {
			return cards, err
		}
		cardWidth = helpers.Max(cardWidth, float64(footerImage.Bounds().Dx()))
		footerCard = render.NewImageContent(render.Style{Width: cardWidth, Height: float64(footerImage.Bounds().Dy())}, footerImage)
	}

	// Header card
	if headerCard, headerCardExists := newHeaderCard(player, options); headerCardExists {
		headerImage, err := headerCard.Render()
		if err != nil {
			return cards, err
		}
		cardWidth = helpers.Max(cardWidth, float64(headerImage.Bounds().Dx()))
		cards = append(cards, render.NewImageContent(render.Style{Width: cardWidth, Height: float64(headerImage.Bounds().Dy())}, headerImage))
	}

	// Player Title card
	cards = append(cards, shared.NewPlayerTitleCard(shared.DefaultPlayerTitleStyle(titleCardStyle(cardWidth)), player.Stats.Account.Nickname, player.Stats.Clan.Tag, player.Subscriptions))

	// Overview Card
	{
		var overviewCardBlocks []render.Block
		for _, column := range player.Cards.Overview.Blocks {
			columnBlock, err := statsBlocksToColumnBlock(getOverviewStyle(overviewColumnWidth), column)
			if err != nil {
				return nil, err
			}
			overviewCardBlocks = append(overviewCardBlocks, columnBlock)
		}
		cards = append(cards, render.NewBlocksContent(overviewCardStyle(cardWidth), overviewCardBlocks...))
	}

	// Highlights
	for _, card := range player.Cards.Highlights {
		cards = append(cards, newHighlightCard(highlightCardStyle(defaultCardStyle(cardWidth)), card))
	}

	// Add footer
	cards = append(cards, footerCard)
	return cards, nil
}

func newHeaderCard(player PlayerData, options RenderOptions) (render.Block, bool) {
	var cards []render.Block

	var addPromoText = true
	for _, sub := range player.Subscriptions {
		switch sub.Type {
		case models.SubscriptionTypePro, models.SubscriptionTypePlus, models.SubscriptionTypeDeveloper:
			addPromoText = false
		}
		if !addPromoText {
			break
		}
	}

	if addPromoText && options.PromoText != nil {
		// Users without a subscription get promo text
		var textBlocks []render.Block
		for _, text := range options.PromoText {
			textBlocks = append(textBlocks, render.NewTextContent(render.Style{Font: &render.FontMedium, FontColor: render.TextPrimary}, text))
		}
		cards = append(cards, render.NewBlocksContent(render.Style{
			Direction:  render.DirectionVertical,
			AlignItems: render.AlignItemsCenter,
		},
			textBlocks...,
		))
	}

	// User Subscription Badge and promo text
	if badges, _ := badges.SubscriptionsBadges(player.Subscriptions); len(badges) > 0 {
		cards = append(cards, render.NewBlocksContent(render.Style{Direction: render.DirectionHorizontal, AlignItems: render.AlignItemsCenter, Gap: 10},
			badges...,
		))
	}

	if len(cards) < 1 {
		return render.Block{}, false
	}

	return render.NewBlocksContent(render.Style{Direction: render.DirectionVertical, AlignItems: render.AlignItemsCenter, JustifyContent: render.JustifyContentCenter, Gap: 10}, cards...), true
}

func newHighlightCard(style highlightStyle, card period.VehicleCard) render.Block {
	titleBlock :=
		render.NewBlocksContent(render.Style{
			Direction: render.DirectionVertical,
		},
			render.NewTextContent(style.cardTitle, card.Meta),
			render.NewTextContent(style.tankName, card.Title),
		)

	var contentRow []render.Block
	for _, block := range card.Blocks {
		contentRow = append(contentRow, render.NewBlocksContent(render.Style{Direction: render.DirectionVertical, AlignItems: render.AlignItemsCenter},
			render.NewTextContent(style.blockValue, block.Data.String),
			render.NewTextContent(style.blockLabel, block.Label),
		))
	}

	return render.NewBlocksContent(style.container, titleBlock, render.NewBlocksContent(render.Style{
		Gap: style.container.Gap,
	}, contentRow...))
}
