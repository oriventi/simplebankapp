package token

import (
	"testing"
	"time"

	"github.com/oriventi/simplebank/util"
	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner()
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, createErr := maker.CreateToken(username, duration)
	require.NoError(t, createErr)
	require.NotEmpty(t, token)

	payload, verifyErr := maker.VerifyToken(token)
	require.NoError(t, verifyErr)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredPasetoMaker(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	token, err := maker.CreateToken(util.RandomOwner(), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, verifyErr := maker.VerifyToken(token)
	require.Error(t, verifyErr)
	require.EqualError(t, verifyErr, ErrExpiredToken.Error())
	require.Nil(t, payload)
}
