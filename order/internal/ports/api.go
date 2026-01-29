package ports

//essa api representa só uma assinatura, um contrato, uma promessa
import (
	"errors"

	"github.com/pauloabaia/microservices/order/internal/application/core/domain"
)

var (
	ErrTooManyItems    = errors.New("order cannot have more than 50 items in total bro, try again with less itens, you are not that rich")
	ErrProductNotFound = errors.New("product was not found in inventory, dumb")
)

type APIPort interface {
	PlaceOrder(order domain.Order) (domain.Order, error)
}
