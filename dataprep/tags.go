package dataprep

import "errors"

type Tag string

const (
	// Global
	TagWN8      Tag = "wn8"
	TagFrags    Tag = "frags"
	TagBattles  Tag = "battles"
	TagWinrate  Tag = "winrate"
	TagAccuracy Tag = "accuracy"

	// Session Specific
	TagAvgDamage   Tag = "avg_damage"
	TagDamageRatio Tag = "damage_ratio"

	// Replay Specific
	TagDamageDealt            Tag = "damage_dealt"
	TagDamageTaken            Tag = "damage_taken"
	TagDamageBlocked          Tag = "blocked"
	TagDamageAssisted         Tag = "assisted"
	TagDamageAssistedCombined Tag = "assisted_combined"
)

func ParseTags(tags ...string) ([]Tag, error) {
	var parsed []Tag
	for _, tag := range tags {
		if tag == "" {
			continue
		}
		switch tag {
		case string(TagWN8):
			parsed = append(parsed, TagWN8)
		case string(TagFrags):
			parsed = append(parsed, TagFrags)
		case string(TagBattles):
			parsed = append(parsed, TagBattles)
		case string(TagWinrate):
			parsed = append(parsed, TagWinrate)
		case string(TagAccuracy):
			parsed = append(parsed, TagAccuracy)

		case string(TagAvgDamage):
			parsed = append(parsed, TagAvgDamage)
		case string(TagDamageRatio):
			parsed = append(parsed, TagDamageRatio)

		case string(TagDamageDealt):
			parsed = append(parsed, TagDamageDealt)
		case string(TagDamageTaken):
			parsed = append(parsed, TagDamageTaken)
		case string(TagDamageBlocked):
			parsed = append(parsed, TagDamageBlocked)
		case string(TagDamageAssisted):
			parsed = append(parsed, TagDamageAssisted)
		case string(TagDamageAssistedCombined):
			parsed = append(parsed, TagDamageAssistedCombined)
		default:
			return nil, errors.New("invalid preset" + tag)
		}
	}

	if len(parsed) == 0 {
		return nil, errors.New("no valid presets")
	}
	return parsed, nil
}
