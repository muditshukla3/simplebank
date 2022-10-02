package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/muditshukla3/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomTestAccount(t *testing.T) (Account, CreateAccountParams, User, error) {
	user := createRandomUser(t)

	arg := CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RandomAmount(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)

	return account, arg, user, err
}
func TestCreateAccount(t *testing.T) {
	account, inputArg, user, err := createRandomTestAccount(t)

	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, inputArg.Owner, account.Owner)
	require.Equal(t, inputArg.Balance, account.Balance)
	require.Equal(t, inputArg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
	// cleanup account
	cleanUpAccount(t, int(account.ID))
	// cleanup user
	cleanUpUser(t, user.Username)
}

func TestGetAccount(t *testing.T) {
	account1, _, user, _ := createRandomTestAccount(t)
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account2)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.Equal(t, account1.ID, account2.ID)

	// cleanup account
	cleanUpAccount(t, int(account1.ID))
	// cleanup user
	cleanUpUser(t, user.Username)
}

func TestDeleteAccount(t *testing.T) {
	account1, _, user, err := createRandomTestAccount(t)
	require.NoError(t, err)
	err = testQueries.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	account, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account)

	// cleanup user
	cleanUpUser(t, user.Username)
}

func TestListAccounts(t *testing.T) {

	lastAccount, _, user, err := createRandomTestAccount(t)
	require.NoError(t, err)
	arg := ListAccountsParams{
		Owner:  lastAccount.Owner,
		Limit:  5,
		Offset: 0,
	}
	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, accounts)

	for _, account := range accounts {
		require.NotEmpty(t, account)
		require.Equal(t, account.Owner, lastAccount.Owner)
	}

	// cleanup account
	cleanUpAccount(t, int(lastAccount.ID))
	// cleanup user
	cleanUpUser(t, user.Username)
}

func cleanUpAccount(t *testing.T, accountId int) {
	err := testQueries.DeleteAccount(context.Background(), int64(accountId))
	require.NoError(t, err)
}
