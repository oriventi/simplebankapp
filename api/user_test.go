package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	mockdb "github.com/oriventi/simplebank/db/mock"
	db "github.com/oriventi/simplebank/db/sqlc"
	"github.com/oriventi/simplebank/util"
	"github.com/stretchr/testify/require"
)

type eqCreateUserArgs struct {
	args     db.CreateUserParams
	password string
}

func (e *eqCreateUserArgs) Matches(x interface{}) bool {
	args := x.(db.CreateUserParams)
	if err := util.CheckPassword(e.password, args.HashedPassword); err != nil {
		return false
	}
	if args.Email != e.args.Email {
		return false
	}
	if args.FullName != e.args.FullName {
		return false
	}
	if args.Username != e.args.Username {
		return false
	}
	return true
}

func (e *eqCreateUserArgs) String() string {
	return fmt.Sprintf("matches %s", e.args.Username)
}

func eqUserArgs(args db.CreateUserParams, password string) gomock.Matcher {
	return &eqCreateUserArgs{args: args, password: password}
}

func TestCreateUser(t *testing.T) {
	hashed, hashErr := util.HashPassword("secret")
	require.NoError(t, hashErr)

	args := db.CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hashed,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}
	reqArgs := createUserRequest{
		Username: args.Username,
		Password: "secret",
		FullName: args.FullName,
		Email:    args.Email,
	}

	user := db.User{
		Username:          args.Username,
		HashedPassword:    hashed,
		FullName:          args.FullName,
		Email:             args.Email,
		PasswortChangedAt: time.Now(),
		CreatedAt:         time.Now(),
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)
	server := newTestServer(t, store)

	store.EXPECT().CreateUser(gomock.Any(), eqUserArgs(args, "secret")).
		Times(1).
		Return(user, nil)

	recorder := httptest.NewRecorder()
	jsonBytes, jsonErr := json.Marshal(reqArgs)
	require.NoError(t, jsonErr)
	req, reqErr := http.NewRequest(http.MethodPost, "/users", bytes.NewReader(jsonBytes))
	require.NoError(t, reqErr)

	server.router.ServeHTTP(recorder, req)
	require.Equal(t, http.StatusOK, recorder.Code)
}
