package stats

import wg "github.com/cufee/am-wg-proxy-next/v2/types"

type AccountWithClan struct {
	wg.ExtendedAccount
	wg.ClanMember
}
