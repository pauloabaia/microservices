package shipping

import (
	"context"
	"log"

	shippingpb "github.com/pauloabaia/microservices-proto/golang/shipping"
	"github.com/pauloabaia/microservices/order/internal/application/core/domain"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Adapter struct {
	shipping shippingpb.ShippingClient
}

func NewAdapter(shippingServiceUrl string) (*Adapter, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.Dial(shippingServiceUrl, opts...)
	if err != nil {
		log.Printf("Failed to dial shipping service: %v", err)
		return nil, err
	}
	client := shippingpb.NewShippingClient(conn)
	return &Adapter{shipping: client}, nil
}

func (a *Adapter) CreateShipping(order *domain.Order) (int32, error) {
	// Converte order items para shipping items
	var items []*shippingpb.ShippingItem
	for _, item := range order.OrderItems {
		items = append(items, &shippingpb.ShippingItem{
			ProductCode: item.ProductCode,
			Quantity:    item.Quantity,
		})
	}

	// Chama o serviço Shipping
	response, err := a.shipping.Create(context.Background(), &shippingpb.CreateShippingRequest{
		OrderId: order.ID,
		Items:   items,
	})
	if err != nil {
		return 0, err
	}

	return response.DeliveryDays, nil
}
