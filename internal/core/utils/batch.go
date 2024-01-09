package utils

func BatchAccountIDs(accountIDs []int, batchSize int) [][]int {
	batches := len(accountIDs) / batchSize
	if len(accountIDs)%100 != 0 {
		batches++
	}

	batchedAccountIDs := make([][]int, batches)
	for i := 0; i < batches; i++ {
		start := i * batchSize
		end := start + batchSize
		if end > len(accountIDs) {
			end = len(accountIDs)
		}

		batchedAccountIDs[i] = accountIDs[start:end]
	}

	return batchedAccountIDs
}
