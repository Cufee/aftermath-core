package replay

import (
	"fmt"

	"github.com/cufee/aftermath-core/internal/core/stats"
	"github.com/cufee/aftermath-core/internal/logic/external"
	"github.com/cufee/aftermath-core/internal/logic/render"
)

type blockPreset struct {
	id    string
	width float64
}

var (
	blockPresetWN8                      = blockPreset{"wn8", 75}
	blockPresetDamageDealt              = blockPreset{"damageDealt", 75}
	blockPresetDamageBlocked            = blockPreset{"damageBlocked", 75}
	blockPresetDamageAssisted           = blockPreset{"damageAssisted", 75}
	blockPresetDamageAssistedAndBlocked = blockPreset{"damageAssistedAndBlocked", 100}
	blockPresetKills                    = blockPreset{"kills", 30}
)

func (preset blockPreset) renderPresetBlock(player *external.Player) render.Block {
	var value int
	switch preset {
	case blockPresetWN8:
		value = player.Performance.WN8(nil)
	case blockPresetDamageDealt:
		value = player.Performance.DamageDealt
	case blockPresetDamageBlocked:
		value = player.Performance.DamageBlocked
	case blockPresetDamageAssisted:
		value = player.Performance.DamageAssisted
	case blockPresetDamageAssistedAndBlocked:
		value = player.Performance.DamageAssisted + player.Performance.DamageBlocked
	case blockPresetKills:
		value = player.Performance.Frags
	}

	text := "-"
	if stats.ValueValid(value) {
		text = fmt.Sprintf("%d", value)
	}

	return render.NewBlocksContent(render.Style{Width: preset.width, JustifyContent: render.JustifyContentCenter}, render.NewTextContent(render.Style{
		Font:      &render.FontLarge,
		FontColor: render.TextPrimary,
	}, text))
}
