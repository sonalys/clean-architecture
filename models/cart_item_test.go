package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	item := CartItem{
		Product:  "",
		Quantity: 1,
	}

	assert.NotNil(t, item.Validate(), "Name should be invalid")
	item.Product = "Pikachu"
	assert.Nil(t, item.Validate(), "Name should be valid")
	item.Quantity = 0
	assert.NotNil(t, item.Validate(), "Quantity should be positive")
}
