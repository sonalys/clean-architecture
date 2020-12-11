package repository

import (
	"errors"

	"github.com/spoonrocker/cart-go-sonalys/models"
	"github.com/spoonrocker/cart-go-sonalys/pkg/persistence"
)

type (
	CartRepository interface {
		New() (*models.Cart, error)
		GetByID(uint64) (*models.Cart, error)
		Update(*models.Cart) (*models.Cart, error)
		RemoveItem(uint64, uint64) error
		AddItem(uint64, *models.CartItem) error
		Delete(*models.Cart) error
	}

	cartRepository struct {
		db persistence.Persistence
	}
)

func CreateCartRepository(p persistence.Persistence) CartRepository {
	return &cartRepository{
		db: p,
	}
}

func (c *cartRepository) New() (*models.Cart, error) {
	cart := models.Cart{
		Items: []models.CartItem{},
	}
	err := c.db.Create(&cart)
	return &cart, err
}

func (c *cartRepository) GetByID(id uint64) (*models.Cart, error) {
	cart := models.Cart{
		ID: id,
	}
	err := c.db.LoadRelationships().Find(&cart)
	return &cart, err
}

func (c *cartRepository) RemoveItem(cartID, cartItemID uint64) error {
	cart, err := c.GetByID(cartID)
	if err != nil || cart == nil {
		return errors.New("cart not found")
	}

	if cart.GetItem(cartItemID) == nil {
		return errors.New("cart does not contain item")
	}

	err = c.db.Delete(&models.CartItem{ID: cartItemID})
	if err != nil {
		return errors.New("error removing cart item")
	}

	return nil
}

func (c *cartRepository) AddItem(cartID uint64, item *models.CartItem) error {
	cart, err := c.GetByID(cartID)
	if err != nil || cart == nil {
		return errors.New("cart not found")
	}

	item.CartID = cart.ID

	err = c.db.Create(item)
	if err != nil {
		return errors.New("error adding cart item")
	}

	return nil
}

func (c *cartRepository) Update(cart *models.Cart) (*models.Cart, error) {
	err := c.db.LoadRelationships().Update(cart)
	return cart, err
}

func (c *cartRepository) Delete(cart *models.Cart) error {
	return c.db.Delete(cart)
}
