package stats

import wg "github.com/cufee/am-wg-proxy-next/types"

type AccountWithClan struct {
	wg.ExtendedAccount
	wg.ClanMember
}
