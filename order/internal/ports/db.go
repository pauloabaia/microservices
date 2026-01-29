package ports

import "github.com/pauloabaia/microservices/order/internal/application/core/domain"

type DBPort interface {
	Get(id string) (domain.Order, error)
	Save(order *domain.Order) error
	ProductExists(productCode string) bool
}
