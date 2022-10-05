package api

import (
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/muditshukla3/simplebank/db/sqlc"
	"github.com/muditshukla3/simplebank/util"
	"github.com/stretchr/testify/require"
)

func NewTestServer(t *testing.T, store db.Store) *Server {

	config := util.Config{
		TokenSymmetricKey:   "UcRefYQrNjcOdpstFsBNFq2yOz9gxThc",
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(config, store)
	require.NoError(t, err)
	return server

}
func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
