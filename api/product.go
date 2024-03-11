package api

import (
	"net/http"

	db "github.com/erlanggawulung/shopifyx/db/sqlc"
	"github.com/gin-gonic/gin"
)

type addProductRequest struct {
	Name           string   `json:"name" binding:"required,min=5,max=60"`
	Price          int      `json:"price" binding:"required,min=0"`
	ImageURL       string   `json:"imageUrl" binding:"required,url"`
	Stock          int      `json:"stock" binding:"required,min=0"`
	Condition      string   `json:"condition" binding:"required,oneof=new second"`
	Tags           []string `json:"tags" binding:"required,min=0"`
	IsPurchaseable bool     `json:"isPurchaseable" binding:"omitempty"`
}

func (server *Server) addProduct(ctx *gin.Context) {
	var req addProductRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	tags := ""
	for i, tag := range req.Tags {
		if i == len(req.Tags)-1 {
			tags += tag
		} else {
			tags += tag + ","
		}
	}
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

	ctx.JSON(http.StatusOK, "")
}