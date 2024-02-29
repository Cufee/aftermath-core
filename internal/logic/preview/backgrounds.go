package preview

import (
	"fmt"
	"image"
	"image/color"
	"sync"

	"github.com/cufee/aftermath-core/internal/core/database"
	"github.com/cufee/aftermath-core/internal/core/utils"
	"github.com/cufee/aftermath-core/internal/logic/content"
	"github.com/cufee/aftermath-core/internal/logic/preview/mock"
	"github.com/cufee/aftermath-core/internal/logic/render"
	"github.com/cufee/aftermath-core/internal/logic/render/session"
	"github.com/cufee/aftermath-core/types"
	wg "github.com/cufee/am-wg-proxy-next/types"
	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
)

var numberBlockWidth = 100

func CurrentBackgroundsPreview() (*types.RenderPreviewResponse, error) {
	data, err := database.GetAppConfiguration[[]string]("backgroundImagesSelection")
	if err != nil {
		return nil, err
	}
	image, err := RenderBackgroundPreview("Your Awesome Nickname", "", data.Value)
	if err != nil {
		return nil, err
	}

	var preview types.RenderPreviewResponse
	preview.Options = data.Value
	preview.Image, err = utils.EncodeImage(image)
	if err != nil {
		return nil, err
	}

	return &preview, nil
}

func RenderBackgroundPreview(nickname, clanTag string, options []string) (image.Image, error) {
	var player session.PlayerData
	player.Cards = mock.PreviewStatsCards
	player.Account = &wg.Account{
		Nickname: nickname,
	}
	if clanTag != "" {
		player.Clan = &wg.Clan{
			Tag: clanTag,
		}
	}

	statsImage, err := renderCardsPreview(player)
	if err != nil {
		return nil, err
	}
	frameWidth, frameHeight := statsImage.Bounds().Dx()+numberBlockWidth, statsImage.Bounds().Dy()*len(options)

	var mux sync.Mutex
	var wait sync.WaitGroup
	previews := make([]utils.DataWithError[image.Image], len(options))
	for i, option := range options {
		wait.Add(1)
		go func(i int, option string) {
			defer wait.Done()

			background, _, err := content.LoadRemoteImage(option)
			if err != nil {
				mux.Lock()
				previews[i] = utils.DataWithError[image.Image]{Err: err}
				mux.Unlock()
				return
			}

			resized := imaging.Fill(background, frameWidth, statsImage.Bounds().Dy(), imaging.Center, imaging.Linear)
			mux.Lock()
			previews[i] = utils.DataWithError[image.Image]{Data: resized}
			mux.Unlock()
		}(i, option)
	}
	wait.Wait()

	frameCtx := gg.NewContext(frameWidth, frameHeight)
	frameCtx.DrawRoundedRectangle(0, 0, float64(frameWidth), float64(frameHeight), 30)
	frameCtx.Clip()

	var lastY int
	for _, preview := range previews {
		if preview.Err != nil {
			return nil, preview.Err
		}
		frameCtx.DrawImage(preview.Data, 0, lastY)
		lastY += preview.Data.Bounds().Dy()
	}

	frameCtx.DrawImage(imaging.Blur(frameCtx.Image(), 10-float64(len(options)-1)), 0, 0)

	lastY = 0
	for i := range previews {
		img, err := statsWithNumber(statsImage, i+1)
		if err != nil {
			return nil, err
		}
		frameCtx.DrawImage(img, 0, lastY)
		lastY += statsImage.Bounds().Dy()
	}

	return frameCtx.Image(), nil
}

func renderCardsPreview(player session.PlayerData) (image.Image, error) {
	renderOptions := session.RenderOptions{
		CardStyle: session.DefaultCardStyle(nil),
	}

	cards, err := session.RenderStatsImage(player, renderOptions)
	if err != nil {
		return nil, err
	}
	return cards, nil
}

func statsWithNumber(statsImage image.Image, number int) (image.Image, error) {
	textBlock := render.NewBlocksContent(render.Style{
		BackgroundColor: color.RGBA{30, 106, 195, 180},
		JustifyContent:  render.JustifyContentCenter,
		AlignItems:      render.AlignItemsCenter,
		BorderRadius:    20,
		Width:           float64(numberBlockWidth - 20),
		Height:          float64(numberBlockWidth - 20),
	}, render.NewTextContent(render.Style{Font: &render.Font2XL, FontColor: render.TextPrimary}, fmt.Sprintf("%d", number)))
	statsBlock := render.NewImageContent(render.Style{Width: float64(statsImage.Bounds().Dx()), Height: float64(statsImage.Bounds().Dy())}, statsImage)
	containerStyle := render.Style{
		JustifyContent: render.JustifyContentStart,
		Direction:      render.DirectionHorizontal,
		AlignItems:     render.AlignItemsCenter,
		Width:          float64(statsImage.Bounds().Dx() + numberBlockWidth),
	}

	if number%2 != 0 {
		containerStyle.JustifyContent = render.JustifyContentEnd
		b := render.NewBlocksContent(containerStyle, textBlock, statsBlock)
		return b.Render()
	} else {
		b := render.NewBlocksContent(containerStyle, statsBlock, textBlock)
		return b.Render()
	}
}
