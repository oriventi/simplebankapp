package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateTransfer(t *testing.T) {

	var amountToTransfer int64 = 20

	acc1 := createTestAccount(t)
	acc2 := createTestAccount(t)
	args := CreateTransferParams{
		FromAccountID: acc1.ID,
		ToAccountID:   acc2.ID,
		Amount:        amountToTransfer,
	}
	transfer, err := testQueries.CreateTransfer(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, transfer.ID)
	require.Equal(t, amountToTransfer, transfer.Amount)
	require.Equal(t, acc1.ID, transfer.FromAccountID)
	require.Equal(t, acc2.ID, transfer.ToAccountID)

	testQueries.DeleteAccount(context.Background(), acc1.ID)
	testQueries.DeleteAccount(context.Background(), acc2.ID)
}

func TestGetTransfer(t *testing.T) {

	var amountToTransfer int64 = 20

	acc1 := createTestAccount(t)
	acc2 := createTestAccount(t)

	args := CreateTransferParams{
		FromAccountID: acc1.ID,
		ToAccountID:   acc2.ID,
		Amount:        amountToTransfer,
	}
	expectedTransfer, err := testQueries.CreateTransfer(context.Background(), args)
	require.NoError(t, err)

	actualTransfer, getErr := testQueries.GetTransfer(context.Background(), expectedTransfer.ID)
	require.NoError(t, getErr)
	require.Equal(t, expectedTransfer, actualTransfer)

	testQueries.DeleteTransfer(context.Background(), actualTransfer.ID)
	testQueries.DeleteAccount(context.Background(), acc1.ID)
	testQueries.DeleteAccount(context.Background(), acc2.ID)
}

func TestListTransfers(t *testing.T) {
	var amountToTransfer int64 = 20
	runs := 5
	fromAccs := make([]Account, runs)
	toAccs := make([]Account, runs)
	transfers := make([]Transfer, runs)

	//CREATE EVERYTHING
	for i := 0; i < runs; i++ {

		fromAccs[i] = createTestAccount(t)
		toAccs[i] = createTestAccount(t)

		args := CreateTransferParams{
			FromAccountID: fromAccs[i].ID,
			ToAccountID:   toAccs[i].ID,
			Amount:        amountToTransfer,
		}
		transfers[i], _ = testQueries.CreateTransfer(context.Background(), args)
	}

	//TEST EVERYTHING
	actualTransfers, err := testQueries.ListTransfers(context.Background(), ListTransfersParams{
		Limit:  int32(runs),
		Offset: 0,
	})
	require.NoError(t, err)
	for i := 0; i < runs; i++ {
		require.NotEmpty(t, actualTransfers[i].ID)
		require.NotEmpty(t, actualTransfers[i].FromAccountID)
		require.NotEmpty(t, actualTransfers[i].ToAccountID)
		require.NotEmpty(t, actualTransfers[i].Amount)
		require.NotEmpty(t, actualTransfers[i].CreatedAt)
	}

	//DELETE EVERYTHING
	for i := 0; i < runs; i++ {
		testQueries.DeleteTransfer(context.Background(), transfers[i].ID)
		testQueries.DeleteAccount(context.Background(), fromAccs[i].ID)
		testQueries.DeleteAccount(context.Background(), toAccs[i].ID)
	}
}
