package main

import (
	"encoding/json"
	"image/png"
	"os"
	"testing"

	"github.com/cufee/aftermath-core/internal/core/database"
	"github.com/cufee/aftermath-core/internal/core/utils"
	"github.com/cufee/aftermath-core/internal/logic/external"
	"github.com/cufee/aftermath-core/internal/logic/render"
	"github.com/cufee/aftermath-core/internal/logic/render/assets"
	"github.com/cufee/aftermath-core/internal/logic/render/replay"
)

func TestFullReplayRenderPipeline(t *testing.T) {
	file, err := os.ReadFile("render_replay_test.json")
	if err != nil {
		panic(err)
	}

	var replayData external.Replay
	err = json.Unmarshal(file, &replayData)
	if err != nil {
		t.Fatal(err)
	}

	err = database.Connect(utils.MustGetEnv("DATABASE_URL"))
	if err != nil {
		t.Fatal(err)
	}

	data := replay.ReplayData{
		Replay: &replayData,
	}

	image, err := replay.RenderReplayImage(data)
	if err != nil {
		t.Fatal(err)
	}

	bgImage, _ := assets.GetImage("images/backgrounds/default")
	img := render.AddBackground(image, bgImage, render.Style{Blur: 10, BorderRadius: 30})

	f, err := os.Create("test-render.png")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	err = png.Encode(f, img)
	if err != nil {
		t.Fatal(err)
	}
}
