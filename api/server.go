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

	// Version 1 routes
	v1 := router.Group("/v1")
	{
		v1.POST("/user/register", server.registerUser)
		v1.POST("/user/login", server.loginUser)
		v1.POST("/product", authMiddleware(server.tokenMaker), server.postProduct)
		v1.PATCH("/product/:id", authMiddleware(server.tokenMaker), server.patchProduct)
		v1.DELETE("/product/:id", authMiddleware(server.tokenMaker), server.deleteProduct)
	}

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
