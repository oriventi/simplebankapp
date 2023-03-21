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
	fromAccs := make(chan Account)
	toAccs := make(chan Account)
	actualFromAccs := make(chan Account)
	actualToAccs := make(chan Account)
	for i := 0; i < n; i++ {
		go func(j int) {
			actualFromAcc, _ := store.CreateRandomAccount()
			actualToAcc, _ := store.CreateRandomAccount()
			fromAcc, toAcc, err := store.execTransfer(context.Background(), TransferTxParams{
				FromAccountID: actualFromAcc.ID,
				ToAccountID:   actualToAcc.ID,
				Amount:        amount,
			})
			errs <- err
			fromAccs <- fromAcc
			actualFromAccs <- actualFromAcc
			actualToAccs <- actualToAcc
			toAccs <- toAcc
		}(i)
	}
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		newFromAcc := <-fromAccs
		actualFromAcc := <-actualFromAccs
		actualToAcc := <-actualToAccs
		require.Equal(t, actualFromAcc.ID, newFromAcc.ID)
		require.Equal(t, actualFromAcc.Balance-amount, newFromAcc.Balance)

		newToAcc := <-toAccs
		require.Equal(t, actualToAcc.ID, newToAcc.ID)
		require.Equal(t, actualToAcc.Balance+amount, newToAcc.Balance)
	}
}
