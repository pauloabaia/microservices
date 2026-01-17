package ports

import (
	"errors"

	"github.com/pauloabaia/microservices/order/internal/application/core/domain"
)

var (
	ErrTooManyItems = errors.New("order cannot have more than 50 items in total bro, try again with less itens, you are not that rich")
)

type APIPort interface {
	PlaceOrder(order domain.Order) (domain.Order, error)
}
