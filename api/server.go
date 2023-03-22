package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/oriventi/simplebank/db/sqlc"
)

// server serves http requests for banking service
type Server struct {
	store  *db.Store
	router *gin.Engine
}

// creates a new httpServer and sets up routing
func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	//add routes
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts", server.listAccounts)
	router.GET("/accounts/:id", server.getAccount)

	server.router = router
	return server
}

func (server *Server) Start(addr string) error {
	return server.router.Run(addr)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
