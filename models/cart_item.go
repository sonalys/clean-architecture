package models

import "errors"

type CartItem struct {
	ID       uint64 `json:"id" gorm:"primaryKey"`
	CartID   uint64 `json:"cart_id"`
	Product  string `json:"product"`
	Quantity uint   `json:"quantity"`
}

func (c *CartItem) Validate() error {
	if c.Product == "" {
		return errors.New("product name cannot be empty")
	}

	if c.Quantity < 1 {
		return errors.New("quantity must be positive")
	}
	return nil
}
