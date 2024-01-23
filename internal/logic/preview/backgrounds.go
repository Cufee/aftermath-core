package preview

import (
	"image"
	"sync"

	"github.com/cufee/aftermath-core/internal/core/database"
	"github.com/cufee/aftermath-core/internal/core/utils"
	"github.com/cufee/aftermath-core/internal/logic/content"
	"github.com/cufee/aftermath-core/internal/logic/preview/mock"
	"github.com/cufee/aftermath-core/internal/logic/render/session"
	"github.com/cufee/aftermath-core/types"
	wg "github.com/cufee/am-wg-proxy-next/types"
	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
)

func CurrentBackgroundsPreview() (*types.RenderPreviewResponse, error) {
	data, err := database.GetAppConfiguration[[]string]("backgroundImagesSelection")
	if err != nil {
		return nil, err
	}
	image, err := RenderBackgroundPreview("Pick Your Style!", "", data.Value)
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

	statsImage, err := renderCardsPreview(player, nil)
	if err != nil {
		return nil, err
	}

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

			resized := imaging.Fill(background, statsImage.Bounds().Dx()/len(options), statsImage.Bounds().Dy(), imaging.Center, imaging.Linear)
			mux.Lock()
			previews[i] = utils.DataWithError[image.Image]{Data: resized}
			mux.Unlock()
		}(i, option)
	}
	wait.Wait()

	lastX := 0
	backgroundWidth := statsImage.Bounds().Dx() / len(options)
	frameCtx := gg.NewContext(statsImage.Bounds().Dx(), statsImage.Bounds().Dy())
	frameCtx.DrawRoundedRectangle(0, 0, float64(statsImage.Bounds().Dx()), float64(statsImage.Bounds().Dy()), 30)
	frameCtx.Clip()

	for _, preview := range previews {
		if preview.Err != nil {
			return nil, preview.Err
		}
		frameCtx.DrawImage(preview.Data, lastX, 0)
		lastX += backgroundWidth
	}
	frameCtx.DrawImage(imaging.Blur(frameCtx.Image(), 10-float64(len(options)-1)), 0, 0)
	frameCtx.DrawImage(statsImage, 0, 0)

	return frameCtx.Image(), nil
}

func renderCardsPreview(player session.PlayerData, background image.Image) (image.Image, error) {
	renderOptions := session.RenderOptions{
		CardStyle: session.DefaultCardStyle(nil),
	}

	cards, err := session.RenderStatsImage(player, renderOptions)
	if err != nil {
		return nil, err
	}
	return cards, nil
}
