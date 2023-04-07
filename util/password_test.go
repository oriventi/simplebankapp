package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHashPassword(t *testing.T) {
	password := "secret"
	hashed, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashed)
	hashed2, err2 := HashPassword(password)
	require.NoError(t, err2)
	require.NotEqual(t, hashed, hashed2)
}

func TestCheckHashed(t *testing.T) {
	password := "secreto"
	hashed, err := HashPassword(password)
	require.NoError(t, err)
	require.NoError(t, CheckPassword(password, hashed))
}
