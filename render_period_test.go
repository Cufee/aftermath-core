package main

import (
	"image/png"
	"os"
	"testing"

	dataprep "github.com/cufee/aftermath-core/dataprep/period"
	"github.com/cufee/aftermath-core/internal/core/database"
	"github.com/cufee/aftermath-core/internal/core/localization"
	"github.com/cufee/aftermath-core/internal/core/utils"
	core "github.com/cufee/aftermath-core/internal/logic/render"
	"github.com/cufee/aftermath-core/internal/logic/render/assets"
	render "github.com/cufee/aftermath-core/internal/logic/render/period"
	"github.com/cufee/aftermath-core/internal/logic/stats/period"
	"golang.org/x/text/language"
)

func TestFullPeriodRenderPipeline(t *testing.T) {
	err := database.Connect(utils.MustGetEnv("DATABASE_URL"))
	if err != nil {
		t.Fatal(err)
	}

	stats, err := period.GetPlayerStats(521493973, 0)
	if err != nil {
		t.Fatal(err)
	}

	var vehicleIDs []int
	for _, vehicle := range stats.Vehicles {
		vehicleIDs = append(vehicleIDs, vehicle.VehicleID)
	}
	vehiclesGlossary, err := database.GetGlossaryVehicles(vehicleIDs...)
	if err != nil {
		t.Fatal(err)
	}

	cards, err := dataprep.SnapshotToSession(dataprep.ExportInput{
		Stats:           stats,
		VehicleGlossary: vehiclesGlossary,
	}, dataprep.ExportOptions{
		Locale:        language.English,
		LocalePrinter: localization.GetPrinter(language.English),

		Blocks:     dataprep.DefaultBlocks,
		Highlights: dataprep.DefaultHighlights,
	})
	if err != nil {
		t.Fatal(err)
	}

	image, err := render.RenderImage(render.PlayerData{
		Stats: stats,
		Cards: cards,
		// Subscriptions: []models.UserSubscription{{Type: models.SubscriptionTypeServerModerator}, {Type: models.SubscriptionTypeServerBooster}, {Type: models.SubscriptionTypePro}, {Type: models.SubscriptionTypeContentTranslator}},
		// Subscriptions: []models.UserSubscription{{Type: models.SubscriptionTypeServerModerator}, {Type: models.SubscriptionTypeServerBooster}, {Type: models.SubscriptionTypeContentTranslator}},
	}, render.RenderOptions{})
	if err != nil {
		t.Fatal(err)
	}

	bgImage, _ := assets.GetImage("images/backgrounds/default")
	img := core.AddBackground(image, bgImage, core.Style{Blur: 10, BorderRadius: 30})

	f, err := os.Create("test-period.png")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	err = png.Encode(f, img)
	if err != nil {
		t.Fatal(err)
	}
}
