package db

import (
	"context"
	"testing"

	"github.com/oriventi/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createTestAccount(t *testing.T) Account {
	userArgs := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: util.RandomString(10),
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}
	user := createTestUser(t, userArgs)

	accArgs := CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RandomBalance(),
		Currency: util.RandomCurrency(),
	}
	acc, err := testQueries.CreateAccount(context.Background(), accArgs)
	require.NoError(t, err)
	require.Equal(t, accArgs.Balance, acc.Balance)
	require.Equal(t, accArgs.Currency, acc.Currency)
	require.Equal(t, accArgs.Owner, acc.Owner)
	return acc
}

func TestCreateAccount(t *testing.T) {
	createTestAccount(t)
}

func TestAddAccountBalance(t *testing.T) {
	n := 5

	for i := 0; i < n; i++ {
		amount := int64(15 + i*3) //something random
		acc := createTestAccount(t)
		newAcc, err := testQueries.AddAccountBalance(context.Background(), AddAccountBalanceParams{
			ID:     acc.ID,
			Amount: amount,
		})
		require.NoError(t, err)
		require.Equal(t, acc.ID, newAcc.ID)
		require.Equal(t, acc.Balance+amount, newAcc.Balance)
	}
}

func TestUpdateAccountBalance(t *testing.T) {
	n := 5

	for i := 0; i < n; i++ {
		newBalance := int64(15 + i*3) //something random
		acc := createTestAccount(t)
		newAcc, err := testQueries.UpdateAccount(context.Background(), UpdateAccountParams{
			ID:      acc.ID,
			Balance: newBalance,
		})
		require.NoError(t, err)
		require.Equal(t, acc.ID, newAcc.ID)
		require.Equal(t, newBalance, newAcc.Balance)
		testQueries.DeleteAccount(context.Background(), acc.ID)
	}
}

func TestDeleteAccount(t *testing.T) {
	n := 5

	for i := 0; i < n; i++ {
		acc := createTestAccount(t)
		foundAcc, err := testQueries.GetAccount(context.Background(), acc.ID)
		require.Equal(t, acc, foundAcc)
		require.NoError(t, err)

		_, err2 := testQueries.DeleteAccount(context.Background(), acc.ID)
		require.NoError(t, err2)

		nullAcc, err3 := testQueries.GetAccount(context.Background(), foundAcc.ID)
		require.Error(t, err3)
		require.Empty(t, nullAcc.ID)

	}
}

func TestGetAccount(t *testing.T) {
	n := 5
	for i := 0; i < n; i++ {
		acc := createTestAccount(t)
		foundAcc, err := testQueries.GetAccount(context.Background(), acc.ID)
		require.NoError(t, err)
		require.Equal(t, acc.ID, foundAcc.ID)
		require.Equal(t, acc.Balance, foundAcc.Balance)
		require.Equal(t, acc.CreatedAt, foundAcc.CreatedAt)
		require.Equal(t, acc.Currency, foundAcc.Currency)
		require.Equal(t, acc.Owner, foundAcc.Owner)

		testQueries.DeleteAccount(context.Background(), acc.ID)
	}
}

func TestListAccounts(t *testing.T) {
	n := 5
	for i := 0; i < n; i++ {
		accs := make([]Account, n)
		for i := 0; i < n; i++ {
			accs[i] = createTestAccount(t)
		}

		foundAccs1, err1 := testQueries.ListAccounts(context.Background(), ListAccountsParams{
			Owner:  accs[n-1].Owner,
			Limit:  int32(n),
			Offset: 0,
		})
		require.NoError(t, err1)
		for i := range foundAccs1 {
			//Test all accounts
			require.NotEmpty(t, foundAccs1[i].ID)
			require.NotEmpty(t, foundAccs1[i].CreatedAt)
			require.NotEmpty(t, foundAccs1[i].Currency)
			require.NotEmpty(t, foundAccs1[i].Owner)

		}
		//Test with different limit and offset
		foundAccs2, err2 := testQueries.ListAccounts(context.Background(), ListAccountsParams{
			Limit:  int32(n / 2),
			Offset: 0,
		})
		require.NoError(t, err2)
		for i := range foundAccs2 {
			require.NotEmpty(t, foundAccs2[i].ID)
			require.NotEmpty(t, foundAccs2[i].Balance)
			require.NotEmpty(t, foundAccs2[i].CreatedAt)
			require.NotEmpty(t, foundAccs2[i].Currency)
			require.NotEmpty(t, foundAccs2[i].Owner)
		}

		for i := 0; i < n; i++ {
			//Delete all accounts
			testQueries.DeleteAccount(context.Background(), accs[i].ID)
		}
	}
}

func TestLockAccount(t *testing.T) {
	n := 5
	for i := 0; i < n; i++ {
		acc := createTestAccount(t)
		tx, txErr := conn.Begin()
		require.NoError(t, txErr)
		q := New(tx)
		foundAcc, err := q.LockAccountEntry(context.Background(), acc.ID)
		require.NoError(t, err)
		require.Equal(t, acc, foundAcc)
		require.Equal(t, acc.Balance, foundAcc.Balance)
		require.Equal(t, acc.CreatedAt, foundAcc.CreatedAt)
		require.Equal(t, acc.Currency, foundAcc.Currency)
		require.Equal(t, acc.Owner, foundAcc.Owner)
		commitErr := tx.Commit()
		require.NoError(t, commitErr)

		testQueries.DeleteAccount(context.Background(), acc.ID)
	}
}
