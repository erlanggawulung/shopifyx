package api

import (
	"fmt"

	db "github.com/erlanggawulung/shopifyx/db/sqlc"
	"github.com/erlanggawulung/shopifyx/token"
	"github.com/erlanggawulung/shopifyx/util"
	"github.com/gin-gonic/gin"
	_ "github.com/gin-gonic/gin/binding"
)

type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewJWTMaker(config.JWTSecret)
	if err != nil {
		return nil, fmt.Errorf("can not create token maker: %v", err)
	}
	server := &Server{config: config, store: store, tokenMaker: tokenMaker}

	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/user/register", server.registerUser)
	router.POST("/user/login", server.loginUser)

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
