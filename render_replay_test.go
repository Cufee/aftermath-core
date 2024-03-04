package main

import (
	"bytes"
	"image/png"
	"os"
	"testing"

	"github.com/cufee/aftermath-core/dataprep"
	replays "github.com/cufee/aftermath-core/dataprep/replay"
	"github.com/cufee/aftermath-core/internal/core/database"
	"github.com/cufee/aftermath-core/internal/core/utils"

	"github.com/cufee/aftermath-core/internal/logic/render"
	"github.com/cufee/aftermath-core/internal/logic/render/assets"
	"github.com/cufee/aftermath-core/internal/logic/render/replay"
	parse "github.com/cufee/aftermath-core/internal/logic/replay"
)

func TestFullReplayRenderPipeline(t *testing.T) {
	file, err := os.ReadFile("render_replay_test.wotbreplay")
	if err != nil {
		panic(err)
	}

	unpacked, err := parse.Unpack(bytes.NewReader(file), int64(len(file)))
	if err != nil {
		t.Fatal(err)
	}
	replayData := parse.Prettify(unpacked.BattleResult, unpacked.Meta)

	err = database.Connect(utils.MustGetEnv("DATABASE_URL"))
	if err != nil {
		t.Fatal(err)
	}

	var vehicles []int
	for _, player := range append(replayData.Teams.Allies, replayData.Teams.Enemies...) {
		vehicles = append(vehicles, player.VehicleID)
	}

	averages, err := database.GetVehicleAverages(vehicles...)
	if err != nil {
		t.Fatal(err)
	}

	cards, err := replays.ReplayToCards(replays.ExportInput{
		GlobalVehicleAverages: averages,
		Replay:                replayData,
	}, replays.ExportOptions{
		Blocks: []dataprep.Tag{dataprep.TagWN8, dataprep.TagDamageDealt, dataprep.TagDamageAssistedCombined, dataprep.TagFrags},
	})
	if err != nil {
		t.Fatal(err)
	}

	image, err := replay.RenderReplayImage(replay.ReplayData{Cards: cards, Replay: replayData}, replay.RenderOptions{})
	if err != nil {
		t.Fatal(err)
	}

	bgImage, _ := assets.GetImage("images/backgrounds/default")
	img := render.AddBackground(image, bgImage, render.Style{Blur: 10, BorderRadius: 30})

	f, err := os.Create("test-replay.png")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	err = png.Encode(f, img)
	if err != nil {
		t.Fatal(err)
	}
}
