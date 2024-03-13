package api

import (
	"database/sql"
	"net/http"

	db "github.com/erlanggawulung/shopifyx/db/sqlc"
	"github.com/erlanggawulung/shopifyx/token"
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
	payload, ok := getContextVal(ctx, authorizationPayloadKey)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Value is not a map"})
		return
	}
	userId := payload.UserId
	productArg := db.CreateProductParams{
		Name:           req.Name,
		Price:          int32(req.Price),
		ImageUrl:       req.ImageURL,
		Stock:          int32(req.Stock),
		Condition:      req.Condition,
		Tags:           tags,
		CreatedBy:      userId,
		IsPurchaseable: req.IsPurchaseable,
	}

	_, err = server.store.CreateProduct(ctx, productArg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

func (server *Server) patchProduct(ctx *gin.Context) {
	// Extract the ID from the URI parameters
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req postProductRequest
	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	tags := sliceToString(req.Tags)
	productArg := db.UpdateProductParams{
		ID:             id,
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

type postProductStockRequest struct {
	Stock int `json:"stock" binding:"required,min=0"`
}

func (server *Server) postProductStock(ctx *gin.Context) {
	// Extract the ID from the URI parameters
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req postProductStockRequest
	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	productArg := db.UpdateProductStockParams{
		ID:    id,
		Stock: int32(req.Stock),
	}

	_, err = server.store.UpdateProductStock(ctx, productArg)
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

type getProductsResponse struct {
	Message string                `json:"message"`
	Data    []productDataResponse `json:"data"`
	Meta    metaResponse          `json:"meta"`
}

type productDataResponse struct {
	ProductID      string   `json:"productId"`
	Name           string   `json:"name"`
	Price          int      `json:"price"`
	ImageURL       string   `json:"imageUrl"`
	Stock          int      `json:"stock"`
	Condition      string   `json:"condition"`
	Tags           []string `json:"tags"`
	IsPurchaseable bool     `json:"isPurchaseable"`
	PurchaseCount  int      `json:"purchaseCount"`
}

type metaResponse struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Total  int `json:"total"`
}

func (server *Server) getProduct(ctx *gin.Context) {

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

func getContextVal(ctx *gin.Context, ctxKey string) (*token.Payload, bool) {
	value, ok := ctx.Get(ctxKey)
	return value.(*token.Payload), ok
}
