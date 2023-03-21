package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(conn)

	n := 5
	const amount = int64(10)
	errs := make(chan error)
	results := make(chan TransferTxResult)
	expectedFromAccs := make(chan Account)
	expectedToAccs := make(chan Account)

	for i := 0; i < n; i++ {
		go func() {
			expectedFromAcc, _ := store.CreateRandomAccount()
			expectedToAcc, _ := store.CreateRandomAccount()
			result, err := store.execTransfer(context.Background(), TransferTxParams{
				FromAccountID: expectedFromAcc.ID,
				ToAccountID:   expectedToAcc.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result
			expectedFromAccs <- expectedFromAcc
			expectedToAccs <- expectedToAcc
		}()
	}
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		expectedFromAcc := <-expectedFromAccs
		expectedToAcc := <-expectedToAccs
		require.Equal(t, expectedFromAcc.ID, result.FromAccount.ID)
		require.Equal(t, expectedFromAcc.Balance-amount, result.FromAccount.Balance)

		require.Equal(t, expectedToAcc.ID, result.ToAccount.ID)
		require.Equal(t, expectedToAcc.Balance+amount, result.ToAccount.Balance)

		require.NotEmpty(t, result.Transfer.ID)
		require.NotEmpty(t, result.Transfer.CreatedAt)
		require.NotEmpty(t, result.Transfer.FromAccountID)
		require.NotEmpty(t, result.Transfer.ToAccountID)

		require.NotEmpty(t, result.FromEntry.CreatedAt)
		require.NotEmpty(t, result.FromEntry.AccountID)
		require.NotEmpty(t, result.FromEntry.ID)

		testQueries.DeleteTransfer(context.Background(), result.Transfer.ID)
		testQueries.DeleteEntry(context.Background(), result.FromEntry.ID)
		testQueries.DeleteEntry(context.Background(), result.ToEntry.ID)

		testQueries.DeleteAccount(context.Background(), result.FromAccount.ID)
		testQueries.DeleteAccount(context.Background(), result.ToAccount.ID)
	}
}
