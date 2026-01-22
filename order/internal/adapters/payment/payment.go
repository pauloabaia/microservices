package payment

import (
	"context"
	"log"
	"time"

	paymentpb "github.com/pauloabaia/microservices-proto/golang/payment"
	"github.com/pauloabaia/microservices/order/internal/application/core/domain"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
)

type Adapter struct {
	payment paymentpb.PaymentClient
}

func NewAdapter(paymentServiceUrl string) (*Adapter, error) {
	var opts []grpc.DialOption
	opts = append(opts,
		grpc.WithUnaryInterceptor(retry.UnaryClientInterceptor( //com o deadline pra forçar o erro
			retry.WithCodes(codes.Unavailable, codes.ResourceExhausted, codes.DeadlineExceeded),
			retry.WithMax(5),
			retry.WithBackoff(retry.BackoffLinear(time.Second)),
		)))
	//opts = append ( opts , grpc.WithInsecure()) modo antigo, depricated
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.Dial(paymentServiceUrl, opts...)
	if err != nil {
		log.Printf("Failed to dial payment service: %v", err)
		return nil, err
	}
	client := paymentpb.NewPaymentClient(conn)
	return &Adapter{payment: client}, nil
}

func (a *Adapter) Charge(order *domain.Order) error {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	_, err := a.payment.Create(ctx, &paymentpb.CreatePaymentRequest{
		UserId:     order.CustomerID,
		OrderId:    order.ID,
		TotalPrice: order.TotalPrice(),
	})

	if err != nil {
		grpcStatus, ok := status.FromError(err)
		if ok && grpcStatus.Code() == codes.DeadlineExceeded {
			log.Printf("Timeout no pagamento: pedido %d excedeu o limite de 10 segundos", order.ID)
		}
	}
	return err

}
