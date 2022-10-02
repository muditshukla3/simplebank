package db

import (
	"context"
	"testing"
	"time"

	"github.com/muditshukla3/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T, account Account) (Entry, error) {
	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    util.RandomAmount(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry, err
}

func TestCreateEntry(t *testing.T) {
	account, _, user, err := createRandomTestAccount(t)
	require.NoError(t, err)
	entry, _ := createRandomEntry(t, account)

	//clean up entry
	cleanUpEntry(t, entry.ID)
	//clean up account
	cleanUpAccount(t, int(account.ID))
	//clean up user
	cleanUpUser(t, user.Username)

}

func TestGetEntry(t *testing.T) {
	account, _, user, _ := createRandomTestAccount(t)
	entry1, _ := createRandomEntry(t, account)
	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, entry1.Amount, entry2.Amount)
	require.WithinDuration(t, entry1.CreatedAt, entry2.CreatedAt, time.Second)

	//clean up entry
	cleanUpEntry(t, entry1.ID)
	//clean up account
	cleanUpAccount(t, int(account.ID))
	//clean up user
	cleanUpUser(t, user.Username)
}

func TestListEntries(t *testing.T) {
	account, _, user, _ := createRandomTestAccount(t)

	entry, _ := createRandomEntry(t, account)

	arg := ListEntriesParams{
		AccountID: account.ID,
		Limit:     5,
		Offset:    0,
	}

	entries, err := testQueries.ListEntries(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entries, 1)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
		require.Equal(t, arg.AccountID, entry.AccountID)
	}

	//clean up entry
	cleanUpEntry(t, entry.ID)
	//clean up account
	cleanUpAccount(t, int(account.ID))
	//clean up user
	cleanUpUser(t, user.Username)
}

func cleanUpEntry(t *testing.T, id int64) {
	err := testQueries.DeleteEntry(context.Background(), id)
	require.NoError(t, err)
}
