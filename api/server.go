package api

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	db "github.com/erlanggawulung/shopifyx/db/sqlc"
	"github.com/erlanggawulung/shopifyx/token"
	"github.com/erlanggawulung/shopifyx/util"
	"github.com/gin-gonic/gin"
	_ "github.com/gin-gonic/gin/binding"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	shopifyxRequestHistogram = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "shopifyx_request",
		Help:    "Histogram of the shopifyx request duration.",
		Buckets: prometheus.LinearBuckets(1, 1, 10), // Adjust bucket sizes as needed
	}, []string{"path", "method", "status"})
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
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Version 1 routes
	v1 := router.Group("/v1")
	{

		NewRoute(v1, "/user/register", http.MethodPost, server.registerUser)
		NewRoute(v1, "/user/login", http.MethodPost, server.loginUser)

		NewRoute(v1, "/product", http.MethodPost, authMiddleware(server.tokenMaker), server.postProduct)
		NewRoute(v1, "/product/:id", http.MethodPatch, authMiddleware(server.tokenMaker), server.patchProduct)
		NewRoute(v1, "/product/:id", http.MethodDelete, authMiddleware(server.tokenMaker), server.deleteProduct)
		NewRoute(v1, "/product/:id/stock", http.MethodPost, authMiddleware(server.tokenMaker), server.postProductStock)

		NewRoute(v1, "/bank/account", http.MethodPost, authMiddleware(server.tokenMaker), server.postBankAccount)
		NewRoute(v1, "/bank/account", http.MethodGet, authMiddleware(server.tokenMaker), server.getBankAccountByUserId)
		NewRoute(v1, "/bank/account/:id", http.MethodDelete, authMiddleware(server.tokenMaker), server.deleteBankAccount)

		NewRoute(v1, "/image", http.MethodPost, authMiddleware(server.tokenMaker), server.postImage)
	}

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func wrapHandlerWithMetrics(path string, method string, handlers ...gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// Execute each handler in sequence
		for _, handler := range handlers {
			handler(c)
			// Check if an error occurred
			if c.IsAborted() {
				break
			}
		}

		// Record the metrics
		duration := time.Since(startTime).Seconds()

		shopifyxRequestHistogram.WithLabelValues(path, method, strconv.Itoa(c.Writer.Status())).Observe(duration)
	}
}

func NewRoute(group *gin.RouterGroup, path, method string, handlers ...gin.HandlerFunc) {
	group.Handle(method, path, wrapHandlerWithMetrics(path, method, handlers...))
}
