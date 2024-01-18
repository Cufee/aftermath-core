package dataprep

import (
	"fmt"

	core "github.com/cufee/aftermath-core/internal/core/stats"
)

func statsToValue(v any) value {
	switch cast := v.(type) {
	case string:
		return value{String: cast, Value: v}
	case float32:
		if int(cast) == core.InvalidValue {
			return value{String: "-", Value: v}
		}
		return value{String: fmt.Sprintf("%.2f", cast), Value: v}
	case float64:
		if int(cast) == core.InvalidValue {
			return value{String: "-", Value: v}
		}
		return value{String: fmt.Sprintf("%.2f%%", cast), Value: v}
	case int:
		if cast == core.InvalidValue {
			return value{String: "-", Value: v}
		}
		return value{String: fmt.Sprint(cast), Value: v}
	default:
		return value{String: "-", Value: v}
	}
}
