package replay

import (
	"errors"
	"fmt"
	"slices"

	"github.com/cufee/aftermath-core/dataprep"
	"github.com/cufee/aftermath-core/internal/core/database"
	"github.com/cufee/aftermath-core/internal/core/localization"
	core "github.com/cufee/aftermath-core/internal/core/stats"
	"github.com/cufee/aftermath-core/internal/core/utils"
	"github.com/cufee/aftermath-core/internal/logic/replay"

	"github.com/rs/zerolog/log"
)

type ExportInput struct {
	Replay                *replay.Replay
	GlobalVehicleAverages map[int]*core.ReducedStatsFrame
}

type ExportOptions struct {
	Blocks []dataprep.Tag
	Locale localization.SupportedLanguage
}

func ReplayToCards(input ExportInput, options ExportOptions) (Cards, error) {
	if input.Replay == nil {
		return Cards{}, errors.New("replay is nil")
	}

	var cards Cards
	printer := localization.GetPrinter(options.Locale)

	var ids []int
	for _, player := range append(input.Replay.Teams.Allies, input.Replay.Teams.Enemies...) {
		ids = append(ids, player.VehicleID)
	}
	vehiclesGlossary, err := database.GetGlossaryVehicles(ids...)
	if err != nil {
		// This is definitely not fatal, but will look ugly
		log.Warn().Err(err).Msg("failed to get vehicles glossary")
	}

	sortTeams(input.Replay.Teams)
	// Allies
	for _, player := range input.Replay.Teams.Allies {
		vehicle := vehiclesGlossary[player.VehicleID]
		vehicle.ID = player.VehicleID
		name := fmt.Sprintf("%s %s", utils.IntToRoman(vehicle.Tier), vehicle.Name(options.Locale))
		cards.Allies = append(cards.Allies, playerToCard(player, name, input.GlobalVehicleAverages[player.VehicleID], options.Blocks, printer))
	}
	// Enemies
	for _, player := range input.Replay.Teams.Enemies {
		vehicle := vehiclesGlossary[player.VehicleID]
		vehicle.ID = player.VehicleID
		name := fmt.Sprintf("%s %s", utils.IntToRoman(vehicle.Tier), vehicle.Name(options.Locale))
		cards.Enemies = append(cards.Enemies, playerToCard(player, name, input.GlobalVehicleAverages[player.VehicleID], options.Blocks, printer))
	}

	return cards, nil
}

func playerToCard(player replay.Player, vehicle string, averages *core.ReducedStatsFrame, blocks []dataprep.Tag, printer localization.LocalePrinter) Card {
	card := Card{
		Type:  dataprep.CardTypeVehicle,
		Meta:  CardMeta{player, blocks},
		Title: vehicle,
	}
	for _, preset := range blocks {
		block := StatsBlock{
			Label: printer("label_" + string(preset)),
			Tag:   preset,
		}
		// Special case
		if preset == dataprep.TagWN8 {
			block.Value = dataprep.Value{Value: float64(player.Performance.WN8(averages)), String: fmt.Sprintf("%d", player.Performance.WN8(averages))}
		} else {
			block.Value = presetToValue(player, preset)
		}
		card.Blocks = append(card.Blocks, block)
	}
	return card
}

func presetToValue(player replay.Player, preset dataprep.Tag) dataprep.Value {
	switch preset {
	case dataprep.TagDamageDealt:
		return dataprep.StatsToValue(player.Performance.DamageDealt)
	case dataprep.TagDamageBlocked:
		if player.Performance.DamageBlocked < 1 {
			return dataprep.Value{Value: -1, String: "-"}
		}
		return dataprep.StatsToValue(player.Performance.DamageBlocked)
	case dataprep.TagDamageAssisted:
		if player.Performance.DamageAssisted < 1 {
			return dataprep.Value{Value: -1, String: "-"}
		}
		return dataprep.StatsToValue(player.Performance.DamageAssisted)
	case dataprep.TagDamageAssistedCombined:
		if player.Performance.DamageAssisted+player.Performance.DamageBlocked < 1 {
			return dataprep.Value{Value: -1, String: "-"}
		}
		return dataprep.StatsToValue(player.Performance.DamageAssisted + player.Performance.DamageBlocked)
	case dataprep.TagFrags:
		if player.Performance.Frags < 1 {
			return dataprep.Value{Value: -1, String: "-"}
		}
		return dataprep.StatsToValue(player.Performance.Frags)
	default:
		return dataprep.Value{Value: -1, String: "-"}
	}
}

func sortTeams(teams replay.Teams) {
	sortPlayers(teams.Allies)
	sortPlayers(teams.Enemies)
}

func sortPlayers(players []replay.Player) {
	slices.SortFunc(players, func(j, i replay.Player) int {
		return (i.Performance.DamageDealt + i.Performance.DamageAssisted + i.Performance.DamageBlocked) - (j.Performance.DamageDealt + j.Performance.DamageAssisted + j.Performance.DamageBlocked)
	})
}
