package api

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	db "github.com/erlanggawulung/shopifyx/db/sqlc"
	"github.com/erlanggawulung/shopifyx/util"
	"github.com/gin-gonic/gin"
)

const (
	tokenDuration = 15 * time.Minute
)

type registerUserRequest struct {
	Username string `json:"username" binding:"required,min=5,max=15"`
	Name     string `json:"name" binding:"required,min=5,max=50"`
	Password string `json:"password" binding:"required,min=5,max=15"`
}

type userResponse struct {
	Message string   `json:"message"`
	Data    userData `json:"data"`
}

type userData struct {
	Username    string `json:"username"`
	Name        string `json:"name"`
	AccessToken string `json:"accessToken"`
}

type loginUserRequest struct {
	Username string `json:"username" binding:"required,min=5,max=15"`
	Password string `json:"password" binding:"required,min=5,max=15"`
}

func generateUserResponse(user db.User, message, accessToken string) userResponse {
	return userResponse{
		Message: message,
		Data: userData{
			Username:    user.Username,
			Name:        user.Name,
			AccessToken: accessToken,
		},
	}
}

func (server *Server) registerUser(ctx *gin.Context) {
	var req registerUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Username: req.Username,
		Password: hashedPassword,
		Name:     req.Name,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		log.Print(db.ErrorCode(err))
		if db.ErrorCode(err) == db.UniqueViolation {
			ctx.JSON(http.StatusConflict, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	accessToken, _, err := server.tokenMaker.CreateToken(req.Username, tokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	res := generateUserResponse(user, "User registered successfully", accessToken)
	ctx.JSON(http.StatusCreated, res)
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	err = util.CheckPassword(req.Password, user.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	accessToken, _, err := server.tokenMaker.CreateToken(req.Username, tokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	res := generateUserResponse(user, "User logged successfully", accessToken)
	ctx.JSON(http.StatusCreated, res)
}
