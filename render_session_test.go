package main

import (
	"image"
	"image/png"
	"os"
	"testing"
	"time"

	"github.com/cufee/aftermath-core/internal/core/localization"
	"github.com/cufee/aftermath-core/internal/logic/render/assets"
	render "github.com/cufee/aftermath-core/internal/logic/render/session"
	"github.com/cufee/aftermath-core/internal/logic/stats"
	"github.com/fogleman/gg"
)

func TestFullRenderPipeline(t *testing.T) {
	start := time.Now()
	session, err := stats.GetCurrentPlayerSession("na", 1013072123) // 1013072123 1032698345
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("got session in %s", time.Since(start).String())

	averages, err := stats.GetVehicleAverages(session.Diff.Vehicles)
	if err != nil {
		t.Fatal(err)
	}

	opts := stats.SortOptions{
		By:    stats.SortByLastBattle,
		Limit: 5,
	}
	player := render.PlayerData{
		Snapshot: session,
		Averages: averages,
		Vehicles: stats.SortVehicles(session.Diff.Vehicles, averages, opts),
	}

	bgImage, _ := assets.GetImage("images/backgrounds/default")
	options := render.RenderOptions{
		PromoText:              []string{"Aftermath is back!", "amth.one/join"},
		Locale:                 localization.LanguageEN,
		CardStyle:              render.DefaultCardStyle(nil),
		UserSubscriptionHeader: nil,
		ClanSubscriptionHeader: render.ClanSubscriptionPremiumXL,
		BackgroundImage:        bgImage,
	}

	now := time.Now()
	tiers := []*render.SubscriptionHeader{nil}
	// tiers := []*render.SubscriptionHeader{nil, render.UserSubscriptionPlus, render.UserSubscriptionPremium, render.UserSubscriptionPremiumXL, render.UserSubscriptionSupporter}
	var images []image.Image
	for _, tier := range tiers {
		optionsWithTier := options
		optionsWithTier.UserSubscriptionHeader = tier
		img, err := render.RenderStatsImage(player, optionsWithTier)
		if err != nil {
			t.Fatal(err)
		}
		images = append(images, img)
	}
	t.Logf("rendered in %s", time.Since(now).String())

	if len(images) == 0 {
		t.Fatal("no images rendered")
	}

	totalHeight := 0
	for _, img := range images {
		totalHeight += img.Bounds().Dy()
	}
	finalCtx := gg.NewContext(images[0].Bounds().Dx(), totalHeight)
	lastY := 0
	for _, img := range images {
		finalCtx.DrawImage(img, 0, lastY)
		lastY += img.Bounds().Dy()
	}

	f, err := os.Create("test.png")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	err = png.Encode(f, finalCtx.Image())
	if err != nil {
		t.Fatal(err)
	}
}
