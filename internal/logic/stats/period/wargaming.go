package period

import (
	"fmt"
	"sync"

	"github.com/cufee/aftermath-core/internal/core/utils"
	"github.com/cufee/aftermath-core/internal/logic/stats"
	"github.com/cufee/am-wg-proxy-next/v2/remote"
	wg "github.com/cufee/am-wg-proxy-next/v2/types"
	"github.com/rs/zerolog/log"
)

func GetAccountInfo(client *remote.Client, realm string, accountID int) (*stats.AccountWithClan, error) {
	var waitGroup sync.WaitGroup

	accountStr := fmt.Sprintf("%d", accountID)
	var account utils.DataWithError[wg.ExtendedAccount]
	var clan wg.ClanMember

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		accounts, err := client.BulkGetAccountsByID([]string{accountStr}, realm, "nickname", "last_battle_time", "account_id", "created_at")
		if err != nil {
			log.Err(err).Msg("failed to get accounts")
		}
		account.Err = err
		account.Data = accounts[accountStr]
	}()

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		clans, err := client.BulkGetAccountsClans([]string{accountStr}, realm, "clan")
		if err != nil {
			// This is not a critical error, so we don't return it
			log.Err(err).Msg("failed to get accounts clans")
		}
		clan = clans[accountStr]
	}()

	waitGroup.Wait()
	if account.Err != nil {
		return nil, account.Err
	}

	return &stats.AccountWithClan{
		ExtendedAccount: account.Data,
		ClanMember:      clan,
	}, nil
}
