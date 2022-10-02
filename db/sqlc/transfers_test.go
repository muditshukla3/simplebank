package db

import (
	"context"
	"testing"
	"time"

	"github.com/muditshukla3/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T, account1, account2 Account) Transfer {
	arg := CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        util.RandomAmount(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)

	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	account1, _, user1, err1 := createRandomTestAccount(t)
	account2, _, user2, err2 := createRandomTestAccount(t)
	require.NoError(t, err1)
	require.NoError(t, err2)
	transfer := createRandomTransfer(t, account1, account2)

	// cleanup transfer
	cleanupTransfer(t, transfer.ID)
	// cleanup account
	cleanUpAccount(t, int(account1.ID))
	cleanUpAccount(t, int(account2.ID))
	// cleanup user
	cleanUpUser(t, user1.Username)
	cleanUpUser(t, user2.Username)
}

func TestGetTransfer(t *testing.T) {
	account1, _, user1, err1 := createRandomTestAccount(t)
	account2, _, user2, err2 := createRandomTestAccount(t)

	require.NoError(t, err1)
	require.NoError(t, err2)

	transfer1 := createRandomTransfer(t, account1, account2)

	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, transfer2)

	require.Equal(t, transfer1.ID, transfer2.ID)
	require.Equal(t, transfer1.FromAccountID, transfer2.FromAccountID)
	require.Equal(t, transfer1.ToAccountID, transfer2.ToAccountID)
	require.Equal(t, transfer1.Amount, transfer2.Amount)
	require.WithinDuration(t, transfer1.CreatedAt, transfer2.CreatedAt, time.Second)

	// cleanup transfer
	cleanupTransfer(t, transfer1.ID)
	// cleanup account
	cleanUpAccount(t, int(account1.ID))
	cleanUpAccount(t, int(account2.ID))
	// cleanup user
	cleanUpUser(t, user1.Username)
	cleanUpUser(t, user2.Username)
}

func TestListTransfer(t *testing.T) {
	account1, _, user1, err1 := createRandomTestAccount(t)
	account2, _, user2, err2 := createRandomTestAccount(t)

	require.NoError(t, err1)
	require.NoError(t, err2)

	transfer := createRandomTransfer(t, account1, account2)

	arg := ListTransfersParams{
		FromAccountID: account1.ID,
		ToAccountID:   account1.ID,
		Limit:         5,
		Offset:        0,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transfers, 1)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
		require.True(t, transfer.FromAccountID == account1.ID || transfer.ToAccountID == account1.ID)
	}

	// cleanup transfer
	cleanupTransfer(t, transfer.ID)
	// cleanup account
	cleanUpAccount(t, int(account1.ID))
	cleanUpAccount(t, int(account2.ID))
	// cleanup user
	cleanUpUser(t, user1.Username)
	cleanUpUser(t, user2.Username)
}

func cleanupTransfer(t *testing.T, id int64) {
	err := testQueries.DeleteTransfer(context.Background(), id)
	require.NoError(t, err)
}
