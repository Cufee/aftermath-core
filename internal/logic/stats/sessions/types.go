package sessions

import (
	"errors"

	core "github.com/cufee/aftermath-core/internal/core/stats"
	"github.com/cufee/aftermath-core/internal/logic/stats"
)

var (
	ErrNoSessionCached = errors.New("no session cached")
)

type Snapshot struct {
	Account  stats.AccountWithClan
	Selected *core.SessionSnapshot // The session that was selected from the database
	Live     *core.SessionSnapshot // The live session
	Diff     *core.SessionSnapshot // The difference between the selected and live sessions
}
