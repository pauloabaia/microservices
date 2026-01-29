package db

import (
	"fmt"

	"github.com/pauloabaia/microservices/shipping/internal/application/core/domain"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Shipping struct {
	gorm.Model
	OrderID      int64
	DeliveryDays int32
}

type Adapter struct {
	db *gorm.DB
}

func (a Adapter) Get(id string) (domain.Shipping, error) {
	var shippingEntity Shipping
	res := a.db.First(&shippingEntity, id)
	shipping := domain.Shipping{
		ID:           int64(shippingEntity.ID),
		OrderID:      shippingEntity.OrderID,
		DeliveryDays: shippingEntity.DeliveryDays,
		CreatedAt:    shippingEntity.CreatedAt.UnixNano(),
	}
	return shipping, res.Error
}

func (a Adapter) Save(shipping *domain.Shipping) error {
	shippingModel := Shipping{
		OrderID:      shipping.OrderID,
		DeliveryDays: shipping.DeliveryDays,
	}
	res := a.db.Create(&shippingModel)
	if res.Error == nil {
		shipping.ID = int64(shippingModel.ID)
	}
	return res.Error
}

func NewAdapter(dataSourceUrl string) (*Adapter, error) {
	db, openErr := gorm.Open(mysql.Open(dataSourceUrl), &gorm.Config{})
	if openErr != nil {
		return nil, fmt.Errorf("db connection error: %v", openErr)
	}
	err := db.AutoMigrate(&Shipping{})
	if err != nil {
		return nil, fmt.Errorf("db migration error: %v", err)
	}
	return &Adapter{db: db}, nil
}