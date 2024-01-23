package main

import (
	"encoding/json"
	"image/png"
	"os"
	"testing"
	"time"

	"github.com/cufee/aftermath-core/dataprep"
	"github.com/cufee/aftermath-core/internal/core/database"
	"github.com/cufee/aftermath-core/internal/core/localization"
	"github.com/cufee/aftermath-core/internal/core/utils"
	"github.com/cufee/aftermath-core/internal/logic/render"
	"github.com/cufee/aftermath-core/internal/logic/render/assets"
	"github.com/cufee/aftermath-core/internal/logic/render/session"
	"github.com/cufee/aftermath-core/internal/logic/stats"
)

func TestFullRenderPipeline(t *testing.T) {
	err := database.Connect(utils.MustGetEnv("DATABASE_URL"))
	if err != nil {
		t.Fatal(err)
	}

	start := time.Now()
	sessionData, err := stats.GetCurrentPlayerSession("na", 1013072123) // 1013072123 1032698345
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("got session in %s", time.Since(start).String())

	averages, err := stats.GetVehicleAverages(sessionData.Diff.Vehicles)
	if err != nil {
		t.Fatal(err)
	}

	opts := stats.SortOptions{
		By:    stats.SortByLastBattle,
		Limit: 5,
	}

	statsBlocks, err := dataprep.SnapshotToSession(dataprep.ExportInput{
		SessionStats:          sessionData.Diff,
		CareerStats:           sessionData.Selected,
		SessionVehicles:       stats.SortVehicles(sessionData.Diff.Vehicles, averages, opts),
		GlobalVehicleAverages: averages,
	}, dataprep.ExportOptions{
		Blocks: dataprep.DefaultBlockPresets,
		Locale: localization.LanguageEN,
	})
	if err != nil {
		t.Fatal(err)
	}

	// session.Account.Nickname = "WWWWWWWWWWWWWWWWWWWWW"
	player := session.PlayerData{
		// Subscriptions: []users.UserSubscription{{Type: users.SubscriptionTypePlus}},
		// Subscriptions: []users.UserSubscription{{Type: users.SubscriptionTypeSupporter}, {Type: users.SubscriptionTypeVerifiedClan}},
		// Subscriptions: []users.UserSubscription{{Type: users.SubscriptionTypeSupporter}, {Type: users.SubscriptionTypeProClan}},
		Account: &sessionData.Account.Account,
		// Clan:    &session.Account.Clan,
		Cards: statsBlocks,
	}

	bgImage, _ := assets.GetImage("images/backgrounds/default")
	options := session.RenderOptions{
		PromoText: []string{"Aftermath is back!", "amth.one/join  |  amth.one/invite"},
		CardStyle: session.DefaultCardStyle(nil),
		// BackgroundImage: bgImage,
	}

	now := time.Now()
	cards, err := session.RenderStatsImage(player, options)
	if err != nil {
		t.Fatal(err)
	}

	data, err := json.MarshalIndent(player.Cards, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	os.WriteFile("test.json", data, 0644)

	img := render.AddBackground(cards, bgImage, render.Style{Blur: 10, BorderRadius: 30, BackgroundColor: render.DiscordBackgroundColor})

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
