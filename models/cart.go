package models

import "errors"

type Cart struct {
	ID    uint64     `json:"id" gorm:"primaryKey"`
	Items []CartItem `json:"items" gorm:"foreignKey:CartID"`
}

func (c *Cart) AddItem(item CartItem) {
	c.Items = append(c.Items, item)
}

func (c *Cart) RemoveItem(id uint64) error {
	for i, cart := range c.Items {
		if cart.ID == id {
			c.Items = append(c.Items[:i], c.Items[i+1:]...)
			return nil
		}
	}
	return errors.New("the item does not exist in the cart")
}

func (c *Cart) GetItem(id uint64) *CartItem {
	for _, cartItem := range c.Items {
		if cartItem.ID == id {
			return &cartItem
		}
	}
	return nil
}
