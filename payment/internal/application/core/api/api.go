package api

import (
	"context"
	"log"
	"time"

	"github.com/huseyinbabal/microservices/payment/internal/application/core/domain"
	"github.com/huseyinbabal/microservices/payment/internal/ports"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Application struct {
	db ports.DBPort
}

func NewApplication(db ports.DBPort) *Application {
	return &Application{
		db: db,
	}
}

func (a Application) Charge(ctx context.Context, payment domain.Payment) (domain.Payment, error) {
	time.Sleep(2 * time.Second) //testar o timeout 
	if payment.TotalPrice > 1000 {
		return domain.Payment{}, status.Errorf(codes.InvalidArgument, "Payment over 1000 is not allowed, youre poor.")
	}
	err := a.db.Save(ctx, &payment)
	if err != nil {
		return domain.Payment{}, err
	}
	return payment, nil
}
