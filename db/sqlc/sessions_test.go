package db

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/muditshukla3/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomSession(t *testing.T) Session {
	user := createRandomUser(t)
	uuid, err := uuid.NewUUID()
	require.NoError(t, err)
	randomToken := util.RandomString(32)
	arg := CreateSessionParams{
		ID:           uuid,
		Username:     user.Username,
		RefreshToken: randomToken,
		UserAgent:    "",
		ClientIp:     "",
		IsBlocked:    false,
		ExpiresAt:    time.Now().Add(time.Hour * 1),
	}

	session, err := testQueries.CreateSession(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, session)
	require.Equal(t, arg.Username, session.Username)
	require.Equal(t, randomToken, session.RefreshToken)
	return session
}

func TestCreateSession(t *testing.T) {
	session := createRandomSession(t)
	//cleanup
	cleanUpSession(t, session.ID)
}

func TestGetSession(t *testing.T) {
	session := createRandomSession(t)

	gotSession, err := testQueries.GetSession(context.Background(), session.ID)
	require.NoError(t, err)
	require.NotEmpty(t, session)
	require.Equal(t, session.Username, gotSession.Username)
	require.Equal(t, session.RefreshToken, gotSession.RefreshToken)

	//cleanup
	cleanUpSession(t, session.ID)
}

func cleanUpSession(t *testing.T, id uuid.UUID) {
	err := testQueries.DeleteSession(context.Background(), id)
	require.NoError(t, err)
}
