package api

import (
	"github.com/pauloabaia/microservices/shipping/internal/application/core/domain"
	"github.com/pauloabaia/microservices/shipping/internal/ports"
)

type Application struct {
	db ports.DBPort
}

func NewApplication(db ports.DBPort) *Application {
	return &Application{
		db: db,
	}
}

func (a Application) CreateShipping(shipping domain.Shipping) (domain.Shipping, error) {
	// Salva o shipping no banco
	err := a.db.Save(&shipping)
	if err != nil {
		return domain.Shipping{}, err
	}
	return shipping, nil
}
