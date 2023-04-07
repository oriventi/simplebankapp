package db

import (
	"context"
	"testing"

	"github.com/oriventi/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createTestUser(t *testing.T, args CreateUserParams) User {
	user, err := testQueries.CreateUser(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	return user
}

func TestCreateUser(t *testing.T) {
	hashedPassword, hashErr := util.HashPassword(util.RandomString(7))
	require.NoError(t, hashErr)
	args := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}
	user := createTestUser(t, args)
	require.Equal(t, args.Username, user.Username)
	require.Equal(t, args.HashedPassword, user.HashedPassword)
	require.Equal(t, args.FullName, user.FullName)
	require.Equal(t, args.Email, user.Email)
}

func TestGetUser(t *testing.T) {
	hashedPassword, hashErr := util.HashPassword(util.RandomString(7))
	require.NoError(t, hashErr)
	args := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}
	expectedUser := createTestUser(t, args)
	actualUser, err := testQueries.GetUser(context.Background(), expectedUser.Username)
	require.NoError(t, err)
	require.Equal(t, args.Username, actualUser.Username)
	require.Equal(t, args.HashedPassword, actualUser.HashedPassword)
	require.Equal(t, args.FullName, actualUser.FullName)
	require.Equal(t, args.Email, actualUser.Email)

}
