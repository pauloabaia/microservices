package grpc

import (
	"context"

	"github.com/pauloabaia/microservices-proto/golang/shipping"
	"github.com/pauloabaia/microservices/shipping/internal/application/core/domain"
)

func (a Adapter) Create(ctx context.Context, request *shipping.CreateShippingRequest) (*shipping.CreateShippingResponse, error) {
	// Converte proto items para domain items
	var items []domain.ShippingItem
	for _, item := range request.Items {
		items = append(items, domain.ShippingItem{
			ProductCode: item.ProductCode,
			Quantity:    item.Quantity,
		})
	}

	// Cria o shipping (calcula prazo automaticamente)
	newShipping := domain.NewShipping(request.OrderId, items)

	// Salva no banco
	result, err := a.api.CreateShipping(newShipping)
	if err != nil {
		return nil, err
	}

	// Retorna resposta
	return &shipping.CreateShippingResponse{
		ShippingId:   result.ID,
		DeliveryDays: result.DeliveryDays,
	}, nil
}
