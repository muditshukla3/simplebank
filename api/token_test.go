package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	mockdb "github.com/muditshukla3/simplebank/db/mock"
	db "github.com/muditshukla3/simplebank/db/sqlc"
	"github.com/muditshukla3/simplebank/token"
	"github.com/muditshukla3/simplebank/util"
	"github.com/stretchr/testify/require"
)

func generateRandomSession(t *testing.T) db.Session {
	maker, err := token.NewPasetoMaker("UcRefYQrNjcOdpstFsBNFq2yOz9gxThc")
	require.NoError(t, err)

	username := util.RandomOwner()
	duration := time.Hour

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, payload, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	return db.Session{
		ID:           payload.ID,
		Username:     payload.Username,
		IsBlocked:    false,
		ExpiresAt:    expiredAt,
		RefreshToken: token,
	}
}
func TestRenewAccessToken(t *testing.T) {
	session := generateRandomSession(t)
	blockedSession := generateRandomSession(t)
	expiredSession := generateRandomSession(t)
	expiredSession.ExpiresAt = time.Now().AddDate(0, 0, -1)
	blockedSession.IsBlocked = true
	tokenResponse := renewAccessTokenResponse{
		AccessToekn:          session.RefreshToken,
		AccessTokenExpiresAt: session.ExpiresAt,
	}

	maker, _ := token.NewPasetoMaker("UcRefYQrNjcOdpstFsBNFq2yOz9gxThc")
	difftoken, _, _ := maker.CreateToken(util.RandomOwner(), time.Minute)

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"refresh_token": session.RefreshToken,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationType, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetSession(gomock.Any(), gomock.Eq(session.ID)).
					Times(1).
					Return(session, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				//check response
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchSession(t, recorder.Body, tokenResponse)
			},
		},
		{
			name: "InvalidBody",
			body: gin.H{
				"refre_token": session.RefreshToken,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationType, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetSession(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				//check response
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "NoSessionFound",
			body: gin.H{
				"refresh_token": session.RefreshToken,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationType, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetSession(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Session{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				//check response
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "InternalServerError",
			body: gin.H{
				"refresh_token": session.RefreshToken,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationType, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetSession(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Session{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				//check response
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "BlockedSession",
			body: gin.H{
				"refresh_token": session.RefreshToken,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationType, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetSession(gomock.Any(), gomock.Any()).
					Times(1).
					Return(blockedSession, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				//check response
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InvalidUser",
			body: gin.H{
				"refresh_token": session.RefreshToken,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationType, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetSession(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Session{Username: util.RandomOwner()}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				//check response
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "TokenDoesNotMatch",
			body: gin.H{
				"refresh_token": difftoken,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationType, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetSession(gomock.Any(), gomock.Any()).
					Times(1).
					Return(session, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				//check response
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "ExpiredToken",
			body: gin.H{
				"refresh_token": difftoken,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationType, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetSession(gomock.Any(), gomock.Any()).
					Times(1).
					Return(expiredSession, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				//check response
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			//start test server and send request
			server := NewTestServer(t, store)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/token/renew_access"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			tc.setupAuth(t, request, server.tokenMaker)
			require.NoError(t, err)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)

		})
	}
}

func requireBodyMatchSession(t *testing.T, body *bytes.Buffer, response renewAccessTokenResponse) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)
	var renewResponse renewAccessTokenResponse
	err = json.Unmarshal(data, &renewResponse)
	require.NoError(t, err)
	require.NotEmpty(t, renewResponse.AccessToekn)
	require.WithinDuration(t, renewResponse.AccessTokenExpiresAt, time.Now(), time.Minute)
}
