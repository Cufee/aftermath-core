package render

import (
	"fmt"
	"image"
	"sync"

	"github.com/cufee/aftermath-core/dataprep"
	replays "github.com/cufee/aftermath-core/dataprep/replay"
	"github.com/cufee/aftermath-core/internal/core/database"
	"github.com/cufee/aftermath-core/internal/core/database/models"
	"github.com/cufee/aftermath-core/internal/core/server"
	core "github.com/cufee/aftermath-core/internal/core/utils"
	"github.com/cufee/aftermath-core/internal/logic/content"
	renderCore "github.com/cufee/aftermath-core/internal/logic/render"
	"github.com/cufee/aftermath-core/internal/logic/render/assets"
	render "github.com/cufee/aftermath-core/internal/logic/render/replay"
	parse "github.com/cufee/aftermath-core/internal/logic/replay"
	"github.com/cufee/aftermath-core/types"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func ReplayFromPayload(c *fiber.Ctx) error {
	var opts types.ReplayRequestPayload
	err := c.BodyParser(&opts)
	if err != nil {
		return c.Status(400).JSON(server.NewErrorResponseFromError(err, "c.BodyParser"))
	}
	if opts.URL == "" {
		return c.Status(400).JSON(server.NewErrorResponse("url is required", "payload"))
	}

	imageData, err := getEncodedReplayImage(opts)
	if err != nil {
		return c.Status(500).JSON(server.NewErrorResponseFromError(err, "getEncodedSessionImage"))
	}

	return c.JSON(server.NewResponse(imageData))
}

func getEncodedReplayImage(options types.ReplayRequestPayload) (string, error) {
	unpacked, err := parse.UnpackRemote(options.URL)
	if err != nil {
		return "", err
	}
	replay := parse.Prettify(unpacked.BattleResult, unpacked.Meta)

	// Fetch the background image in a separate goroutine
	var wait sync.WaitGroup
	backgroundChan := make(chan image.Image, 1)
	cardsChan := make(chan core.DataWithError[image.Image], 1)

	wait.Add(1)
	go func() {
		defer wait.Done()

		referenceIDs := []string{fmt.Sprint(replay.Protagonist.ID), fmt.Sprint(replay.Protagonist.ClanID)}
		backgrounds, err := database.GetContentByReferenceIDs[string](referenceIDs, models.UserContentTypePersonalBackground, models.UserContentTypeClanBackground)
		if err != nil {
			log.Warn().Err(err).Msg("failed to get backgrounds")
			bgImage, _ := assets.GetImage("images/backgrounds/default")
			backgroundChan <- bgImage
			return
		}

		// We should get personal image over clan image when possible, fallback to default
		for _, id := range referenceIDs {
			for _, c := range backgrounds {
				if c.Data != "" && c.ReferenceID == id {
					image, _, err := content.LoadRemoteImage(c.Data)
					if err == nil && image != nil {
						backgroundChan <- image
						return
					}
				}
			}
		}
		// fallback
		bgImage, _ := assets.GetImage("images/backgrounds/default")
		backgroundChan <- bgImage
	}()

	wait.Add(1)
	go func() {
		defer wait.Done()

		var vehicles []int
		for _, player := range append(replay.Teams.Allies, replay.Teams.Enemies...) {
			vehicles = append(vehicles, player.VehicleID)
		}

		averages, err := database.GetVehicleAverages(vehicles...)
		if err != nil {
			cardsChan <- core.DataWithError[image.Image]{Err: err}
			return
		}

		cards, err := replays.ReplayToCards(replays.ExportInput{
			GlobalVehicleAverages: averages,
			Replay:                replay,
		}, replays.ExportOptions{
			Blocks: []dataprep.Tag{dataprep.TagWN8, dataprep.TagDamageDealt, dataprep.TagDamageAssistedCombined, dataprep.TagFrags},
		})
		if err != nil {
			cardsChan <- core.DataWithError[image.Image]{Err: err}
			return
		}

		img, err := render.RenderReplayImage(render.ReplayData{Cards: cards, Replay: replay})
		cardsChan <- core.DataWithError[image.Image]{Data: img, Err: err}
	}()

	wait.Wait()
	close(cardsChan)
	close(backgroundChan)

	cards := <-cardsChan
	if cards.Err != nil {
		return "", cards.Err
	}

	bgImage := <-backgroundChan
	img := renderCore.AddBackground(cards.Data, bgImage, renderCore.Style{Blur: 10, BorderRadius: 30})
	return core.EncodeImage(img)
}
