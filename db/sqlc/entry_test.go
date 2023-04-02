package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCreateEntry(t *testing.T) {

	var amountInEntry int64 = 15
	acc := createTestAccount(t)

	args := CreateEntryParams{
		AccountID: acc.ID,
		Amount:    int64(amountInEntry),
		CreatedAt: time.Now(),
	}
	entry, err := testQueries.CreateEntry(context.Background(), args)
	require.NoError(t, err)
	require.Equal(t, acc.ID, entry.AccountID)
	require.Equal(t, amountInEntry, entry.Amount)
	require.NotEmpty(t, entry.ID)
	testQueries.DeleteEntry(context.Background(), entry.ID)
	testQueries.DeleteAccount(context.Background(), acc.ID)
}

func TestGetEntry(t *testing.T) {
	var amountInEntry int64 = 15
	acc := createTestAccount(t)

	args := CreateEntryParams{
		AccountID: acc.ID,
		Amount:    int64(amountInEntry),
		CreatedAt: time.Now(),
	}
	expectedEntry, _ := testQueries.CreateEntry(context.Background(), args)
	actualEntry, err := testQueries.GetEntry(context.Background(), expectedEntry.ID)
	require.NoError(t, err)
	require.Equal(t, expectedEntry, actualEntry)

	testQueries.DeleteEntry(context.Background(), expectedEntry.ID)
	testQueries.DeleteAccount(context.Background(), acc.ID)
}

func TestDeleteEntry(t *testing.T) {
	var amountInEntry int64 = 15
	acc := createTestAccount(t)

	args := CreateEntryParams{
		AccountID: acc.ID,
		Amount:    int64(amountInEntry),
		CreatedAt: time.Now(),
	}
	expectedEntry, _ := testQueries.CreateEntry(context.Background(), args)
	_, err := testQueries.DeleteEntry(context.Background(), expectedEntry.ID)

	require.NoError(t, err)

	_, err = testQueries.GetEntry(context.Background(), expectedEntry.ID)

	require.Equal(t, sql.ErrNoRows, err)
}

func TestListEntries(t *testing.T) {
	runs := 5
	expectedEntries := make([]Entry, runs)

	//CREATE EVERYTHING
	for i := 0; i < runs; i++ {
		var amountInEntry int64 = 15
		acc := createTestAccount(t)

		args := CreateEntryParams{
			AccountID: acc.ID,
			Amount:    int64(amountInEntry),
			CreatedAt: time.Now(),
		}
		expectedEntries[i], _ = testQueries.CreateEntry(context.Background(), args)
	}

	//TEST EVERYTHING
	for i := 0; i < runs; i++ {
		actualEntries, err := testQueries.ListEntries(context.Background(), ListEntriesParams{
			Limit:  int32(runs),
			Offset: 0,
		})
		require.NoError(t, err)
		require.Equal(t, len(actualEntries), runs)
		require.NotEmpty(t, actualEntries[i].ID)
		require.NotEmpty(t, actualEntries[i].AccountID)
		require.NotEmpty(t, actualEntries[i].Amount)
	}

	//DELETE EVERYTHING
	for i := 0; i < runs; i++ {
		testQueries.DeleteEntry(context.Background(), expectedEntries[i].ID)
	}
}
