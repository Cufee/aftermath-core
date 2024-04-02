package utils

import "github.com/cufee/am-wg-proxy-next/v2/utils"

func RealmFromAccountID(accountID int) string {
	return utils.RealmFromPlayerID(accountID)
}
