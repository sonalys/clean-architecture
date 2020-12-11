package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddItem(t *testing.T) {
	cart := Cart{
		ID:    1,
		Items: []CartItem{},
	}

	assert.Len(t, cart.Items, 0, "Should have 0 items")
	cart.AddItem(CartItem{ID: 1})
	assert.Len(t, cart.Items, 1, "Should have 1 item")
}

func TestRemoveItem(t *testing.T) {
	cart := Cart{
		ID: 1,
		Items: []CartItem{
			{ID: 1},
			{ID: 2},
			{ID: 3},
			{ID: 4},
			{ID: 5},
		},
	}

	assert.Len(t, cart.Items, 5, "Should have 5 items")
	cart.RemoveItem(1)
	assert.Equal(t, cart.Items, []CartItem{
		{ID: 2},
		{ID: 3},
		{ID: 4},
		{ID: 5},
	}, "Should remove first item")

	cart.RemoveItem(5)
	assert.Equal(t, cart.Items, []CartItem{
		{ID: 2},
		{ID: 3},
		{ID: 4},
	}, "Should remove last item")

	cart.RemoveItem(3)
	assert.Equal(t, cart.Items, []CartItem{
		{ID: 2},
		{ID: 4},
	}, "Should remove middle item")

	cart.RemoveItem(2)

	assert.Nil(t, cart.RemoveItem(4),
		"Should not create error when sucessfuly removing item")

	assert.Len(t, cart.Items, 0, "Should be empty")

	assert.NotNil(t, cart.RemoveItem(99),
		"Should create error when not finding an item")
}

func TestContainsItem(t *testing.T) {
	cart := Cart{
		Items: []CartItem{
			{ID: 1},
			{ID: 2},
			{ID: 3},
		},
	}

	assert.Equal(t, uint64(1), cart.GetItem(1).ID, "Should find first item")
	assert.Equal(t, uint64(2), cart.GetItem(2).ID, "Should find middle item")
	assert.Equal(t, uint64(3), cart.GetItem(3).ID, "Should find last item")
	assert.Nil(t, cart.GetItem(4), "Should not find inexistent item")

	cart.Items = []CartItem{}
	assert.Nil(t, cart.GetItem(4), "Should find nothing on empty list")
}
