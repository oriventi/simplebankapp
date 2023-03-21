package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/oriventi/simplebank/util"
)

// Store provides all functions to execute db queries and transactions
type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a function within a database transaction
func (store *Store) execTx(ctx context.Context, callback func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, &sql.TxOptions{
		ReadOnly: false,
	})

	if err != nil {
		return err
	}

	q := New(tx)
	defer tx.Rollback()

	err2 := callback(q)
	if err2 != nil {
		return err2
	}
	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountID int64
	ToAccountID   int64
	Amount        int64
}

type TransferTxResult struct {
	FromAccount Account
	ToAccount   Account
	FromEntry   Entry
	ToEntry     Entry
	Transfer    Transfer
}

func (store *Store) execTransfer(ctx context.Context, args TransferTxParams) (result TransferTxResult, err error) {
	err = store.execTx(ctx, func(q *Queries) error {

		//create transfer
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: args.FromAccountID,
			ToAccountID:   args.ToAccountID,
			Amount:        args.Amount,
		})
		if err != nil {
			return err
		}

		//create entries
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: args.FromAccountID,
			Amount:    -args.Amount,
			CreatedAt: time.Now(),
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: args.ToAccountID,
			Amount:    args.Amount,
			CreatedAt: time.Now(),
		})
		if err != nil {
			return err
		}

		//lock accounts
		q.LockAccountEntry(ctx, args.FromAccountID)
		q.LockAccountEntry(ctx, args.ToAccountID)

		//update accounts
		result.FromAccount, result.ToAccount, err = addMoneyToTwoAccounts(
			ctx,
			q, args.FromAccountID, -args.Amount,
			args.ToAccountID, args.Amount,
		)
		if err != nil {
			return err
		}

		return nil
	})

	return
}

func (store *Store) CreateRandomAccount() (Account, error) {
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomBalance(),
		Currency: util.RandomCurrency(),
	}
	createdAcc, err := store.CreateAccount(context.Background(), arg)
	return createdAcc, err
}

func addMoneyToTwoAccounts(
	ctx context.Context,
	q *Queries,
	accountID1,
	amount1,
	accountID2,
	amount2 int64,
) (account1 Account, account2 Account, err error) {

	//add money to account 1
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		Amount: amount1,
		ID:     accountID1,
	})
	if err != nil {
		return
	}

	//move money from account 2
	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		Amount: amount2,
		ID:     accountID2,
	})
	return
}
