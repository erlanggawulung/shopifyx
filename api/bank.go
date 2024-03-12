package api

import (
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
	userId := payload.ID
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
	userId := payload.ID

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
