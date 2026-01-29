package ports

import "github.com/pauloabaia/microservices/order/internal/application/core/domain"

type ShippingPort interface {
	CreateShipping(order *domain.Order) (int32, error)
}
