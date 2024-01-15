package session

import (
	"fmt"

	core "github.com/cufee/aftermath-core/internal/core/stats"
)

func statsValueToString(value any) string {
	switch cast := value.(type) {
	case string:
		return cast
	case float64:
		if int(cast) == core.InvalidValue {
			return "-"
		}
		return fmt.Sprintf("%.2f%%", value)
	case int:
		if value == core.InvalidValue {
			return "-"
		}
		return fmt.Sprintf("%d", value)
	default:
		return fmt.Sprint(value)
	}
}
