package api

import (
	"database/sql"
	"log"
	"net/http"

	db "github.com/erlanggawulung/shopifyx/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	successMessage = "success"
)

type postBankAccountRequest struct {
	BankName          string `json:"bankName" binding:"required,min=5,max=15"`
	BankAccountName   string `json:"bankAccountName" binding:"required,min=5,max=15"`
	BankAccountNumber string `json:"bankAccountNumber" binding:"required,min=5,max=15"`
}

func (server *Server) postBankAccount(ctx *gin.Context) {
	var req postBankAccountRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload, ok := getContextVal(ctx, authorizationPayloadKey)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Value is not a map"})
		return
	}
	userId := payload.UserId
	bankArg := db.CreateBankAccountParams{
		BankName:          req.BankName,
		BankAccountName:   req.BankAccountName,
		BankAccountNumber: req.BankAccountNumber,
		UserID:            userId,
	}

	_, err = server.store.CreateBankAccount(ctx, bankArg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

type getBankAccountsByUserIdResponse struct {
	Message string                     `json:"message"`
	Data    []bankAccountsDataResponse `json:"data"`
}

type bankAccountsDataResponse struct {
	BankAccountId     uuid.UUID `json:"bankAccountId"`
	BankName          string    `json:"bankName"`
	BankAccountName   string    `json:"bankAccountName"`
	BankAccountNumber string    `json:"bankAccountNumber"`
}

func (server *Server) getBankAccountByUserId(ctx *gin.Context) {
	payload, ok := getContextVal(ctx, authorizationPayloadKey)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Value is not a map"})
		return
	}
	userId := payload.UserId

	bankAccounts, err := server.store.GetBankAccountsByUserId(ctx, userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	message := successMessage
	bankAccList := generateBankAccountsResponse(bankAccounts)

	ctx.JSON(http.StatusOK, getBankAccountsByUserIdResponse{
		Message: message,
		Data:    bankAccList,
	})
}

func generateBankAccountsResponse(bankAccounts []db.GetBankAccountsByUserIdRow) []bankAccountsDataResponse {
	results := []bankAccountsDataResponse{}
	for _, val := range bankAccounts {
		bankAcc := bankAccountsDataResponse{
			BankAccountId:     val.ID,
			BankName:          val.BankName,
			BankAccountName:   val.BankAccountName,
			BankAccountNumber: val.BankAccountNumber,
		}
		results = append(results, bankAcc)
	}
	return results
}

func (server *Server) deleteBankAccount(ctx *gin.Context) {
	// Extract the ID from the URI parameters
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload, ok := getContextVal(ctx, authorizationPayloadKey)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Value is not a map"})
		return
	}
	userId := payload.UserId

	deleteBankAccountArg := db.DeleteBankAccountParams{
		ID:     id,
		UserID: userId,
	}
	log.Print(deleteBankAccountArg)
	_, err = server.store.DeleteBankAccount(ctx, deleteBankAccountArg)
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

type patchBankAccountRequest struct {
	BankAccountId     string `json:"bankAccountId" binding:"required"`
	BankName          string `json:"bankName" binding:"required,min=5,max=15"`
	BankAccountName   string `json:"bankAccountName" binding:"required,min=5,max=15"`
	BankAccountNumber string `json:"bankAccountNumber" binding:"required,min=5,max=15"`
}

func (server *Server) patchBankAccount(ctx *gin.Context) {
	// Extract the ID from the URI parameters
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req patchBankAccountRequest
	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload, ok := getContextVal(ctx, authorizationPayloadKey)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Value is not a map"})
		return
	}
	userId := payload.UserId

	updateBankAccArg := db.UpdateBankAccountParams{
		ID:                id,
		UserID:            userId,
		BankName:          req.BankName,
		BankAccountName:   req.BankAccountName,
		BankAccountNumber: req.BankAccountNumber,
	}
	_, err = server.store.UpdateBankAccount(ctx, updateBankAccArg)
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
