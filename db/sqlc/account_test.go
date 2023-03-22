package db

import (
	"context"
	"testing"

	"github.com/oriventi/simplebank/util"
	"github.com/stretchr/testify/require"
)

func TestCreateAccount(t *testing.T) {
	n := 5

	for i := 0; i < n; i++ {
		arg := CreateAccountParams{
			Owner:    util.RandomOwner(),
			Balance:  util.RandomBalance(),
			Currency: util.RandomCurrency(),
		}
		acc, err := testQueries.CreateAccount(context.Background(), arg)
		require.NoError(t, err)
		require.NotEmpty(t, acc)
	}
}

func TestAddAccountBalance(t *testing.T) {
	n := 5

	for i := 0; i < n; i++ {
		amount := int64(15 + i*3) //something random
		arg := CreateAccountParams{
			Owner:    util.RandomOwner(),
			Balance:  util.RandomBalance(),
			Currency: util.RandomCurrency(),
		}
		acc, _ := testQueries.CreateAccount(context.Background(), arg)
		newAcc, err := testQueries.AddAccountBalance(context.Background(), AddAccountBalanceParams{
			ID:     acc.ID,
			Amount: amount,
		})
		require.NoError(t, err)
		require.Equal(t, acc.ID, newAcc.ID)
		require.Equal(t, acc.Balance+amount, newAcc.Balance)
		testQueries.DeleteAccount(context.Background(), acc.ID)
	}
}

func TestUpdateAccountBalance(t *testing.T) {
	n := 5

	for i := 0; i < n; i++ {
		newBalance := int64(15 + i*3) //something random
		arg := CreateAccountParams{
			Owner:    util.RandomOwner(),
			Balance:  util.RandomBalance(),
			Currency: util.RandomCurrency(),
		}
		acc, _ := testQueries.CreateAccount(context.Background(), arg)
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
		arg := CreateAccountParams{
			Owner:    util.RandomOwner(),
			Balance:  util.RandomBalance(),
			Currency: util.RandomCurrency(),
		}
		acc, _ := testQueries.CreateAccount(context.Background(), arg)
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
		arg := CreateAccountParams{
			Owner:    util.RandomOwner(),
			Balance:  util.RandomBalance(),
			Currency: util.RandomCurrency(),
		}
		acc, _ := testQueries.CreateAccount(context.Background(), arg)

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
			arg := CreateAccountParams{
				Owner:    util.RandomOwner(),
				Balance:  util.RandomBalance(),
				Currency: util.RandomCurrency(),
			}
			accs[i], _ = testQueries.CreateAccount(context.Background(), arg)
		}

		foundAccs1, err1 := testQueries.ListAccounts(context.Background(), ListAccountsParams{
			Limit:  int32(n),
			Offset: 0,
		})
		require.NoError(t, err1)
		for i := range accs {
			//Test all accounts
			require.NotEmpty(t, foundAccs1[i].ID)
			require.NotEmpty(t, foundAccs1[i].Balance)
			require.NotEmpty(t, foundAccs1[i].CreatedAt)
			require.NotEmpty(t, foundAccs1[i].Currency)
			require.NotEmpty(t, foundAccs1[i].Owner)

		}
		//Test with different limit and offset
		foundAccs2, err2 := testQueries.ListAccounts(context.Background(), ListAccountsParams{
			Limit:  int32(n / 2),
			Offset: 2,
		})
		require.NoError(t, err2)
		for i := 0; i < int(n/2); i++ {
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
		arg := CreateAccountParams{
			Owner:    util.RandomOwner(),
			Balance:  util.RandomBalance(),
			Currency: util.RandomCurrency(),
		}
		acc, _ := testQueries.CreateAccount(context.Background(), arg)
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
