package ports

import "github.com/pauloabaia/microservices/shipping/internal/application/core/domain"

type DBPort interface {
	Save(shipping *domain.Shipping) error
	Get(id string) (domain.Shipping, error)
}
