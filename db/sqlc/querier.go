// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.21.0

package db

import (
	"context"

	"github.com/google/uuid"
)

type Querier interface {
	CreateBankAccount(ctx context.Context, arg CreateBankAccountParams) (BankAccount, error)
	CreateProduct(ctx context.Context, arg CreateProductParams) (Product, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	DeleteProduct(ctx context.Context, id uuid.UUID) (Product, error)
	GetBankAccountsByUserId(ctx context.Context, userID uuid.UUID) ([]GetBankAccountsByUserIdRow, error)
	GetProduct(ctx context.Context, id uuid.UUID) (Product, error)
	GetUser(ctx context.Context, username string) (User, error)
	ListProducts(ctx context.Context, arg ListProductsParams) (Product, error)
	UpdateProduct(ctx context.Context, arg UpdateProductParams) (Product, error)
	UpdateProductStock(ctx context.Context, arg UpdateProductStockParams) (Product, error)
}

var _ Querier = (*Queries)(nil)
