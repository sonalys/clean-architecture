package usecase

import (
	"strconv"

	"github.com/spoonrocker/cart-go-sonalys/models"
	"github.com/spoonrocker/cart-go-sonalys/pkg/http"
	"github.com/spoonrocker/cart-go-sonalys/pkg/persistence"
	"github.com/spoonrocker/cart-go-sonalys/repository"
)

type (
	CartUsecase interface {
		CreateCart(http.HandlerContext) error
		AddItem(http.HandlerContext) error
		RemoveItem(http.HandlerContext) error
		ListItems(ctx http.HandlerContext) error
	}

	cartUsecase struct {
		repository repository.CartRepository
	}
)

const (
	INVALID_CART_ID      = "invalid cart id"
	INVALID_CART_ITEM_ID = "invalid cart item id"
	CART_NOT_FOUND       = "cart not found"
	INVALID_CART_ITEM    = "cart item is not valid"
	FAILED_ADD_ITEM      = "failed to add item to cart"
)

func CreateCartUsecase(httpHandler http.HTTP, db persistence.Persistence) CartUsecase {
	var cartUsecase = cartUsecase{
		repository: repository.CreateCartRepository(db),
	}

	httpHandler.AddRoute("POST", "/carts", cartUsecase.CreateCart)
	httpHandler.AddRoute("POST", "/carts/:cart_id/items", cartUsecase.AddItem)
	httpHandler.AddRoute("GET", "/carts/:cart_id", cartUsecase.ListItems)
	httpHandler.AddRoute("DELETE", "/carts/:cart_id/items/:cart_item_id", cartUsecase.RemoveItem)

	return cartUsecase
}

func (c cartUsecase) CreateCart(ctx http.HandlerContext) error {
	cart, err := c.repository.New()
	if err != nil {
		return err
	}

	return ctx.JSON(200, cart)
}

func (c cartUsecase) ListItems(ctx http.HandlerContext) error {
	cartID, err := strconv.ParseUint(ctx.Param("cart_id"), 0, 64)
	if err != nil {
		return ctx.String(404, INVALID_CART_ID)
	}

	cart, err := c.repository.GetByID(cartID)
	if err != nil || cart == nil {
		return ctx.String(404, CART_NOT_FOUND)
	}

	return ctx.JSON(200, cart)
}

func (c cartUsecase) AddItem(ctx http.HandlerContext) error {
	cartID, err := strconv.ParseUint(ctx.Param("cart_id"), 0, 64)

	if err != nil {
		return ctx.String(400, INVALID_CART_ID)
	}

	cart, err := c.repository.GetByID(cartID)
	if err != nil {
		return ctx.String(404, CART_NOT_FOUND)
	}

	cartItem := new(models.CartItem)
	err = ctx.Bind(cartItem)
	if err != nil {
		return ctx.String(406, INVALID_CART_ITEM)
	}

	if err = cartItem.Validate(); err != nil {
		return ctx.String(406, err.Error())
	}

	err = c.repository.AddItem(cartID, cartItem)
	if err != nil || cart == nil {
		return ctx.String(500, FAILED_ADD_ITEM)
	}

	return ctx.JSON(200, cartItem)
}

func (c cartUsecase) RemoveItem(ctx http.HandlerContext) error {
	cartID, err := strconv.ParseUint(ctx.Param("cart_id"), 0, 64)
	if err != nil {
		return ctx.String(400, INVALID_CART_ID)
	}

	cartItemID, err := strconv.ParseUint(ctx.Param("cart_item_id"), 0, 64)
	if err != nil {
		return ctx.String(400, INVALID_CART_ITEM_ID)
	}

	err = c.repository.RemoveItem(cartID, cartItemID)
	if err != nil {
		return ctx.String(500, err.Error())
	}

	return ctx.JSON(200, struct{}{})
}
