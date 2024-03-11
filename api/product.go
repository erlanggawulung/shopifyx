package api

import (
	"database/sql"
	"net/http"

	db "github.com/erlanggawulung/shopifyx/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type postProductRequest struct {
	Name           string   `json:"name" binding:"required,min=5,max=60"`
	Price          int      `json:"price" binding:"required,min=0"`
	ImageURL       string   `json:"imageUrl" binding:"required,url"`
	Stock          int      `json:"stock" binding:"required,min=0"`
	Condition      string   `json:"condition" binding:"required,oneof=new second"`
	Tags           []string `json:"tags" binding:"required,min=0"`
	IsPurchaseable bool     `json:"isPurchaseable" binding:"omitempty"`
}

func (server *Server) postProduct(ctx *gin.Context) {
	var req postProductRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	tags := sliceToString(req.Tags)
	productArg := db.CreateProductParams{
		Name:           req.Name,
		Price:          int32(req.Price),
		ImageUrl:       req.ImageURL,
		Stock:          int32(req.Stock),
		Condition:      req.Condition,
		Tags:           tags,
		IsPurchaseable: req.IsPurchaseable,
	}

	_, err = server.store.CreateProduct(ctx, productArg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

type patchProductRequest struct {
	ID             uuid.UUID `uri:"id" binding:"required,min=1"`
	Name           string    `json:"name" binding:"required,min=5,max=60"`
	Price          int       `json:"price" binding:"required,min=0"`
	ImageURL       string    `json:"imageUrl" binding:"required,url"`
	Stock          int       `json:"stock" binding:"required,min=0"`
	Condition      string    `json:"condition" binding:"required,oneof=new second"`
	Tags           []string  `json:"tags" binding:"required,min=0"`
	IsPurchaseable bool      `json:"isPurchaseable" binding:"omitempty"`
}

func (server *Server) patchProduct(ctx *gin.Context) {
	// Extract the ID from the URI parameters
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req patchProductRequest
	req.ID = id

	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	tags := sliceToString(req.Tags)
	productArg := db.UpdateProductParams{
		ID:             req.ID,
		Name:           req.Name,
		Price:          int32(req.Price),
		ImageUrl:       req.ImageURL,
		Stock:          int32(req.Stock),
		Condition:      req.Condition,
		Tags:           tags,
		IsPurchaseable: req.IsPurchaseable,
	}

	_, err = server.store.UpdateProduct(ctx, productArg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

func (server *Server) deleteProduct(ctx *gin.Context) {
	// Extract the ID from the URI parameters
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err = server.store.DeleteProduct(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

func sliceToString(slice []string) string {
	results := ""
	for i, tag := range slice {
		if i == len(slice)-1 {
			results += tag
		} else {
			results += tag + ","
		}
	}
	return results
}
