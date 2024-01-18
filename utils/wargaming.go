package utils

import "github.com/cufee/am-wg-proxy-next/client"

func RealmFromAccountID(accountID int) string {
	return client.RealmFromPlayerID(accountID)
}
