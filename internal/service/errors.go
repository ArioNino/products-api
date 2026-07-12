package service

import (
	"errors"
)

var (
	ErrProductNameRequired = errors.New("product name is required")
	ErrProductStockRequired = errors.New("product price is required")
	ErrProductPriceRequired = errors.New("product stock is required")
)