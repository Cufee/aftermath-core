package utils

import "github.com/cufee/am-wg-proxy-next/utils"

func RealmFromAccountID(accountID int) string {
	return utils.RealmFromPlayerID(accountID)
}
