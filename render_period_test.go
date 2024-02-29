package main

import (
	"image/png"
	"os"
	"testing"

	dataprep "github.com/cufee/aftermath-core/dataprep/period"
	"github.com/cufee/aftermath-core/internal/core/database"
	"github.com/cufee/aftermath-core/internal/core/database/models"
	"github.com/cufee/aftermath-core/internal/core/utils"
	render "github.com/cufee/aftermath-core/internal/logic/render/period"
	"github.com/cufee/aftermath-core/internal/logic/stats/period"
)

func TestFullPeriodRenderPipeline(t *testing.T) {
	err := database.Connect(utils.MustGetEnv("DATABASE_URL"))
	if err != nil {
		t.Fatal(err)
	}

	stats, err := period.GetPlayerStats(1013072123, 0)
	if err != nil {
		t.Fatal(err)
	}

	cards, err := dataprep.SnapshotToSession(stats, dataprep.ExportOptions{Blocks: dataprep.DefaultBlocks})
	if err != nil {
		t.Fatal(err)
	}

	image, err := render.RenderImage(render.PlayerData{
		Stats:         stats,
		Cards:         cards,
		Subscriptions: []models.UserSubscription{{Type: models.SubscriptionTypeServerModerator}, {Type: models.SubscriptionTypeServerBooster}, {Type: models.SubscriptionTypePro}, {Type: models.SubscriptionTypeContentTranslator}},
	}, render.RenderOptions{})
	if err != nil {
		t.Fatal(err)
	}

	f, err := os.Create("test-period.png")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	err = png.Encode(f, image)
	if err != nil {
		t.Fatal(err)
	}
}
