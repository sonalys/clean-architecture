package repository

import (
	"errors"
	"testing"

	"github.com/jinzhu/copier"
	"github.com/spoonrocker/cart-go-sonalys/models"
	"github.com/spoonrocker/cart-go-sonalys/pkg/persistence"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	cart := models.Cart{
		ID: 1,
	}

	fakeDB := persistence.MockPersistence{
		MockCreate: func(i interface{}) error {
			copier.Copy(i, cart)
			return nil
		},
	}

	repo := cartRepository{
		db: fakeDB,
	}

	createdCart, err := repo.New()
	assert.Nil(t, err, "Should not have error")
	assert.Equal(t, cart, *createdCart)
}

func FakeDB(cart models.Cart) persistence.Persistence {
	return persistence.MockPersistence{
		MockLoadRelationships: func() persistence.Persistence {
			return persistence.MockPersistence{
				MockFind: func(i interface{}) error {
					if i.(*models.Cart).ID != cart.ID {
						return errors.New("Searching for wrong ID")
					}
					copier.Copy(i, cart)
					return nil
				},
			}
		},
		MockDelete: func(i interface{}) error {
			if i.(*models.CartItem).ID != cart.Items[0].ID {
				return errors.New("Deleting wrong cart item")
			}
			return nil
		},
	}
}
func TestGetByID(t *testing.T) {
	cart := models.Cart{
		ID: 1,
	}

	repo := cartRepository{
		db: FakeDB(cart),
	}

	createdCart, err := repo.GetByID(cart.ID)
	assert.Nil(t, err, "Should not have error")
	assert.Equal(t, cart, *createdCart)
}

func TestRemoveItem(t *testing.T) {
	cart := models.Cart{
		ID: 1,
		Items: []models.CartItem{
			{ID: 123},
		},
	}

	fakeDB := FakeDB(cart)

	repo := cartRepository{
		db: fakeDB,
	}

	err := repo.RemoveItem(cart.ID, cart.Items[0].ID)
	assert.Nil(t, err, "Should not have error")
}
