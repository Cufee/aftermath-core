package main

import (
	"encoding/json"
	"image/png"
	"os"
	"testing"
	"time"

	sessionPrep "github.com/cufee/aftermath-core/dataprep/session"
	"github.com/cufee/aftermath-core/internal/core/database"
	"github.com/cufee/aftermath-core/internal/core/database/models"
	"github.com/cufee/aftermath-core/internal/core/utils"
	"github.com/cufee/aftermath-core/internal/logic/render"
	"github.com/cufee/aftermath-core/internal/logic/render/assets"
	"github.com/cufee/aftermath-core/internal/logic/render/session"
	"github.com/cufee/aftermath-core/internal/logic/stats"
)

func TestFullSessionRenderPipeline(t *testing.T) {
	file, err := os.ReadFile("render_session_test.json")
	if err != nil {
		panic(err)
	}

	var statsCards sessionPrep.Cards
	err = json.Unmarshal(file, &statsCards)
	if err != nil {
		t.Fatal(err)
	}

	err = database.Connect(utils.MustGetEnv("DATABASE_URL"))
	if err != nil {
		t.Fatal(err)
	}

	sessionData, err := stats.GetCurrentPlayerSession("eu", 581650793) // 1013072123 1032698345
	if err != nil {
		t.Fatal(err)
	}

	// session.Account.Nickname = "WWWWWWWWWWWWWWWWWWWWW"
	player := session.PlayerData{
		Subscriptions: []models.UserSubscription{{Type: models.SubscriptionTypeServerModerator}},
		// Subscriptions: []models.UserSubscription{{Type: models.SubscriptionTypeSupporter}, {Type: models.SubscriptionTypeVerifiedClan}},
		// Subscriptions: []models.UserSubscription{{Type: models.SubscriptionTypeServerModerator}, {Type: models.SubscriptionTypeServerBooster}, {Type: models.SubscriptionTypePro}, {Type: models.SubscriptionTypeContentTranslator}},
		Account: &sessionData.Account.Account,
		Session: sessionData,
		// Clan:    &session.Account.Clan,
		Cards: statsCards,
	}

	bgImage, _ := assets.GetImage("images/backgrounds/light")
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

	img := render.AddBackground(cards, bgImage, render.Style{Blur: 10, BorderRadius: 30})
	t.Logf("rendered in %s", time.Since(now).String())

	f, err := os.Create("test-session.png")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	err = png.Encode(f, img)
	if err != nil {
		t.Fatal(err)
	}
}
