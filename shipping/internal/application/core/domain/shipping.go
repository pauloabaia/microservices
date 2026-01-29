package domain

import "time"

type ShippingItem struct {
	ProductCode string
	Quantity    int32
}

type Shipping struct {
	ID           int64
	OrderID      int64
	DeliveryDays int32
	CreatedAt    int64
}

func NewShipping(orderId int64, items []ShippingItem) Shipping {
	return Shipping{
		OrderID:      orderId,
		DeliveryDays: getDeliveryDays(items),
		CreatedAt:    time.Now().Unix(),
	}
}

// Calcula prazo: mínimo 1 dia + 1 dia a cada 5 unidades
func getDeliveryDays(items []ShippingItem) int32 { //el pega os itens 
	var totalQuantity int32
	for _, item := range items {
		totalQuantity += item.Quantity
	}
	days := 1 + (totalQuantity / 5)

	return days
}
