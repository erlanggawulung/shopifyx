// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.21.0

package db

import (
	"time"

	"github.com/google/uuid"
)

type BankAccount struct {
	ID                uuid.UUID `json:"id"`
	UserID            uuid.UUID `json:"user_id"`
	BankName          string    `json:"bank_name"`
	BankAccountName   string    `json:"bank_account_name"`
	BankAccountNumber string    `json:"bank_account_number"`
	CreatedAt         time.Time `json:"created_at"`
}

type Product struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	Price          int32     `json:"price"`
	ImageUrl       string    `json:"image_url"`
	Stock          int32     `json:"stock"`
	Condition      string    `json:"condition"`
	Tags           string    `json:"tags"`
	IsPurchaseable bool      `json:"is_purchaseable"`
	CreatedBy      uuid.UUID `json:"created_by"`
	CreatedAt      time.Time `json:"created_at"`
}

type User struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Name      string    `json:"name"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}
