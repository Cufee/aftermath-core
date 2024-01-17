package main

import (
	"image/png"
	"os"
	"testing"
	"time"

	"github.com/cufee/aftermath-core/internal/core/localization"
	"github.com/cufee/aftermath-core/internal/logic/dataprep"
	"github.com/cufee/aftermath-core/internal/logic/render/assets"
	render "github.com/cufee/aftermath-core/internal/logic/render/session"
	"github.com/cufee/aftermath-core/internal/logic/stats"
	"github.com/cufee/aftermath-core/internal/logic/users"
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
	session.Account.Clan.Tag = "AMTH"

	statsBlocks, err := dataprep.SnapshotToSession(dataprep.ExportInput{
		SessionStats:          session.Diff,
		CareerStats:           session.Selected,
		SessionVehicles:       stats.SortVehicles(session.Diff.Vehicles, averages, opts),
		GlobalVehicleAverages: averages,
	}, dataprep.ExportOptions{
		Blocks: dataprep.DefaultBlockPresets,
		Locale: localization.LanguageEN,
	})
	if err != nil {
		t.Fatal(err)
	}

	player := render.PlayerData{
		// Subscriptions: []users.UserSubscription{{Type: users.SubscriptionTypeProClan}},
		Subscriptions: []users.UserSubscription{{Type: users.SubscriptionTypeSupporter}, {Type: users.SubscriptionTypeProClan}},
		Snapshot:      session,
		Blocks:        &statsBlocks,
	}

	bgImage, _ := assets.GetImage("images/backgrounds/default")
	options := render.RenderOptions{
		PromoText:       []string{"Aftermath is back!", "amth.one/join"},
		Locale:          localization.LanguageEN,
		CardStyle:       render.DefaultCardStyle(nil),
		BackgroundImage: bgImage,
	}

	now := time.Now()
	img, err := render.RenderStatsImage(player, options)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("rendered in %s", time.Since(now).String())

	f, err := os.Create("test.png")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	err = png.Encode(f, img)
	if err != nil {
		t.Fatal(err)
	}
}
