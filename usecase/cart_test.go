package usecase

import (
	"errors"
	"testing"

	"github.com/jinzhu/copier"
	"github.com/spoonrocker/cart-go-sonalys/models"
	"github.com/spoonrocker/cart-go-sonalys/pkg/http"
	"github.com/spoonrocker/cart-go-sonalys/pkg/persistence"
	"github.com/stretchr/testify/assert"
)

func simplifiedErrorFunc(i int, s string) error {
	return errors.New(s)
}

func fakeDB(cart models.Cart) persistence.Persistence {
	return persistence.MockPersistence{
		MockLoadRelationships: func() persistence.Persistence {
			return persistence.MockPersistence{
				MockFind: func(i interface{}) error {
					copier.Copy(i, &cart)
					return nil
				},
			}
		},
		MockDelete: func(i interface{}) error {
			if i.(*models.CartItem).ID != cart.Items[0].ID {
				return errors.New("Wrong cart item id")
			}
			return nil
		},
		MockUpdate: func(i interface{}) error { return nil },
		MockCreate: func(i interface{}) error { return nil },
	}
}

func TestAddItem(t *testing.T) {

	cart := models.Cart{
		ID: 1,
		Items: []models.CartItem{
			{ID: 1},
		},
	}

	persistenceMock := fakeDB(cart)
	cartUsercase := CreateCartUsecase(
		&http.MockHTTP{}, persistenceMock)
	ctxMock := http.MockHandlerContext{
		MockParam:  func(s string) string { return "a" },
		MockString: simplifiedErrorFunc,
		MockJSON:   func(i1 int, i2 interface{}) error { return nil },
	}

	err := cartUsercase.AddItem(ctxMock)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), INVALID_CART_ID, "Should give error about wrong id")

	mockProduct := &models.CartItem{
		ID:      2,
		Product: "",
	}

	ctxMock.MockParam = func(s string) string { return "2" }
	ctxMock.MockBind = func(i interface{}) error {
		copier.Copy(i, mockProduct)
		return nil
	}

	assert.Equal(
		t, cartUsercase.AddItem(ctxMock).Error(), mockProduct.Validate().Error(), "Should validate item name")

	mockProduct.Product = "foo"
	assert.Equal(
		t, cartUsercase.AddItem(ctxMock).Error(), mockProduct.Validate().Error(), "Quantity should be positive")

	mockProduct.Quantity = 1
	assert.Nil(t, cartUsercase.AddItem(ctxMock), "Error should be nil")
}

func TestRemoveItem(t *testing.T) {
	cart := models.Cart{
		ID: 1,
		Items: []models.CartItem{
			{ID: 1},
		},
	}

	persistenceMock := fakeDB(cart)
	cartUsercase := CreateCartUsecase(
		&http.MockHTTP{}, persistenceMock)

	ParamResponseMap := map[string]string{
		"cart_id":      "a",
		"cart_item_id": "b",
	}

	ctxMock := http.MockHandlerContext{
		MockParam:  func(s string) string { return ParamResponseMap[s] },
		MockString: simplifiedErrorFunc,
		MockJSON:   func(i1 int, i2 interface{}) error { return nil },
	}

	err := cartUsercase.RemoveItem(ctxMock)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), INVALID_CART_ID, "Should give error about wrong cart id")

	ParamResponseMap["cart_id"] = "2"

	err = cartUsercase.RemoveItem(ctxMock)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), INVALID_CART_ITEM_ID, "Should give error about wrong cart item id")

	ParamResponseMap["cart_item_id"] = "1"

	assert.Nil(t, cartUsercase.RemoveItem(ctxMock), "Error should be nil")
}

func TestListItem(t *testing.T) {
	cart := models.Cart{
		ID: 1,
		Items: []models.CartItem{
			{ID: 1},
		},
	}

	persistenceMock := fakeDB(cart)
	cartUsercase := CreateCartUsecase(
		&http.MockHTTP{}, persistenceMock)
	ctxMock := http.MockHandlerContext{
		MockParam:  func(s string) string { return "a" },
		MockString: simplifiedErrorFunc,
		MockJSON:   func(i1 int, i2 interface{}) error { return nil },
	}

	err := cartUsercase.ListItems(ctxMock)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), INVALID_CART_ID, "Should give error about wrong cart id")

	ctxMock.MockParam = func(s string) string { return "1" }

	assert.Nil(t, cartUsercase.ListItems(ctxMock), "Error should be nil")
}
