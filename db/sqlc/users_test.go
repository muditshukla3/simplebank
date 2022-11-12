package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/muditshukla3/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	hashedPwd, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)
	arg := CreateUserParams{
		Username: util.RandomOwner(),
		Password: hashedPwd,
		FullName: util.RandomOwner(),
		Email:    util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.Password, user.Password)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)

	require.NotZero(t, user.CreatedAt)

	return user
}
func TestCreateUser(t *testing.T) {
	user := createRandomUser(t)

	//cleanup
	cleanUpUser(t, user.Username)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.Password, user2.Password)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)

	//cleanup
	cleanUpUser(t, user1.Username)
}

func cleanUpUser(t *testing.T, username string) {
	err := testQueries.DeleteUser(context.Background(), username)
	require.NoError(t, err)
}

func TestUserUpdateOnlyFullName(t *testing.T) {
	oldUser := createRandomUser(t)
	newFullname := util.RandomOwner()
	updatedUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		FullName: sql.NullString{
			String: newFullname,
			Valid:  true,
		},
	})
	require.NoError(t, err)
	require.NotEqual(t, oldUser.FullName, updatedUser.FullName)
	require.Equal(t, newFullname, updatedUser.FullName)
	require.Equal(t, oldUser.Email, updatedUser.Email)
	require.Equal(t, oldUser.Password, updatedUser.Password)

	//cleanup
	cleanUpUser(t, oldUser.Username)
}

func TestUserUpdateOnlyEmail(t *testing.T) {
	oldUser := createRandomUser(t)
	newEmail := util.RandomEmail()
	updatedUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		Email: sql.NullString{
			String: newEmail,
			Valid:  true,
		},
	})
	require.NoError(t, err)
	require.NotEqual(t, oldUser.Email, updatedUser.FullName)
	require.Equal(t, newEmail, updatedUser.Email)
	require.Equal(t, oldUser.FullName, updatedUser.FullName)
	require.Equal(t, oldUser.Password, updatedUser.Password)

	//cleanup
	cleanUpUser(t, oldUser.Username)
}

func TestUserUpdateOnlyPassword(t *testing.T) {
	oldUser := createRandomUser(t)
	newPassword := util.RandomString(6)
	newHashedPassword, err := util.HashPassword(newPassword)
	require.NoError(t, err)
	updatedUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		Password: sql.NullString{
			String: newHashedPassword,
			Valid:  true,
		},
	})

	require.NoError(t, err)
	require.NotEqual(t, oldUser.Password, updatedUser.Password)
	require.Equal(t, newHashedPassword, updatedUser.Password)
	require.Equal(t, oldUser.FullName, updatedUser.FullName)
	require.Equal(t, oldUser.Email, updatedUser.Email)
	//cleanup
	cleanUpUser(t, oldUser.Username)
}

func TestUserUpdateAllFields(t *testing.T) {
	oldUser := createRandomUser(t)
	newEmail := util.RandomEmail()
	newFullName := util.RandomOwner()
	newPassword := util.RandomString(6)
	newHashedPassword, err := util.HashPassword(newPassword)
	require.NoError(t, err)
	updatedUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		Email: sql.NullString{
			String: newEmail,
			Valid:  true,
		},
		Password: sql.NullString{
			String: newHashedPassword,
			Valid:  true,
		},
		FullName: sql.NullString{
			String: newFullName,
			Valid:  true,
		},
	})
	require.NoError(t, err)
	require.NotEqual(t, oldUser.Email, updatedUser.FullName)
	require.NotEqual(t, oldUser.Password, updatedUser.Password)
	require.NotEqual(t, oldUser.FullName, updatedUser.FullName)
	require.Equal(t, newEmail, updatedUser.Email)
	require.Equal(t, newFullName, updatedUser.FullName)
	require.Equal(t, newHashedPassword, updatedUser.Password)

	//cleanup
	cleanUpUser(t, oldUser.Username)
}
