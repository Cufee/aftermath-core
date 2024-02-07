package dataprep

import (
	"fmt"

	"github.com/cufee/aftermath-core/internal/core/stats"
)

func StatsToValue(v any) Value {
	switch cast := v.(type) {
	case float32:
		if cast == stats.InvalidValueFloat32 {
			return Value{String: "-", Value: stats.InvalidValueFloat64}
		}
		return Value{String: fmt.Sprintf("%.2f", cast), Value: float64(cast)}
	case float64:
		if cast == stats.InvalidValueFloat64 {
			return Value{String: "-", Value: stats.InvalidValueFloat64}
		}
		return Value{String: fmt.Sprintf("%.2f%%", cast), Value: cast}
	case int:
		if cast == stats.InvalidValueInt {
			return Value{String: "-", Value: stats.InvalidValueFloat64}
		}
		return Value{String: fmt.Sprint(cast), Value: float64(cast)}
	default:
		return Value{String: "-", Value: stats.InvalidValueFloat64}
	}
}
