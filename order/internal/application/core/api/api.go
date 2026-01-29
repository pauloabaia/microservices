package api

//essa api aqui representa lógica de negócio, lá é assinatura e aqui é implementação
import (
	"github.com/pauloabaia/microservices/order/internal/application/core/domain"
	"github.com/pauloabaia/microservices/order/internal/ports"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Application struct {
	db       ports.DBPort
	payment  ports.PaymentPort
	shipping ports.ShippingPort
}

func NewApplication(db ports.DBPort, payment ports.PaymentPort, shipping ports.ShippingPort) *Application {
	return &Application{
		db:       db,
		payment:  payment,
		shipping: shipping,
	}
}

func (a Application) PlaceOrder(order domain.Order) (domain.Order, error) {
	// Validar se todos os produtos existem no estoque
	for _, item := range order.OrderItems {
		if !a.db.ProductExists(item.ProductCode) {
			return domain.Order{}, status.Errorf(codes.InvalidArgument, "%s: %s", ports.ErrProductNotFound.Error(), item.ProductCode)
		}
	}

	// Validar total de itens (soma das quantidades)
	var totalItems int32
	for _, item := range order.OrderItems {
		totalItems += item.Quantity
	}
	if totalItems > 50 {
		return domain.Order{}, status.Errorf(codes.InvalidArgument, "%s", ports.ErrTooManyItems.Error())
	}

	err := a.db.Save(&order) // Salva com status "Pending"
	if err != nil {
		return domain.Order{}, err
	}

	paymentErr := a.payment.Charge(&order)
	if paymentErr != nil {
		order.Status = "Canceled" // Muda o status
		a.db.Save(&order)         // GORM faz UPDATE (pois order.ID já existe)
		return domain.Order{}, paymentErr
	}

	// Payment teve sucesso, chama Shipping
	deliveryDays, shippingErr := a.shipping.CreateShipping(&order)
	if shippingErr != nil {
		// Shipping falhou mas pagamento já foi aprovado
		deliveryDays = 0 // Prazo desconhecido se shipping falhar
	}

	order.DeliveryDays = deliveryDays
	order.Status = "Paid" // Muda o status
	a.db.Save(&order)     // GORM faz UPDATE
	return order, nil
}
