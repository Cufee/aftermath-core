package main

import (
	"image/png"
	"os"
	"testing"
	"time"

	"github.com/cufee/aftermath-core/internal/core/localization"
	"github.com/cufee/aftermath-core/internal/logic/render"
	"github.com/cufee/aftermath-core/internal/logic/stats"
)

func TestFullRenderPipeline(t *testing.T) {
	start := time.Now()
	session, err := stats.GetCurrentPlayerSession("na", 1032698345) // 1013072123 1032698345
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("got session in %s", time.Since(start).String())

	averages, err := stats.GetVehicleAverages(session.Diff.Vehicles)
	if err != nil {
		t.Fatal(err)
	}

	// TODO: sorting options and limits
	opts := stats.SortOptions{
		By:    stats.SortByLastBattle,
		Limit: 7,
	}
	vehicles := stats.SortVehicles(session.Diff.Vehicles, averages, opts)

	start = time.Now()
	img, err := render.RenderStatsImage(session, vehicles, averages, localization.LanguageEN)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("rendered in %s", time.Since(start).String())

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
