package repository

import (
	"errors"

	"github.com/mhmmmdrivaldhi/go-book-api/model"
	"gorm.io/gorm"
)

type OrderRepository interface {
	CreateOrder(order *model.Order) (*model.Order, error)
	FindAll() ([]model.Order, error)
	FindById(id int) (*model.Order, error)
	UpdateOrder(id int, updateOrder *model.Order) (*model.Order, error)
	DeleteOrder(id int) error
}

type orderRepository struct {
	db *gorm.DB
}

func(or *orderRepository) CreateOrder(order *model.Order) (*model.Order, error) {
	err := or.db.Transaction(func(tx *gorm.DB) error  {
		if err := tx.Create(&order).Error; err != nil {
			return err
		}
	
		for  item := range order.OrderItems {
			order.OrderItems[item].OrderID = order.ID
			err := tx.Create(&order.OrderItems[item]).Error
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return order, nil
}

func (or *orderRepository) FindAll() ([]model.Order, error) {
	var orders []model.Order

	err := or.db.Preload("User").Preload("OrderItems").Find(&orders).Error
	if err != nil {
		return nil, err
	}

	if len(orders) == 0 {
		return nil, errors.New("no orders found")
	}

	return orders, nil
}

func (or *orderRepository) FindById(id int) (*model.Order, error) {
	var order model.Order

	err := or.db.Preload("User").Preload("OrderItems").First(&order, id).Error
	if err != nil {
		return nil, err
	}
	
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("order not found")
	}

	return &order, nil
}

func (or *orderRepository) UpdateOrder(id int, updateOrder *model.Order) (*model.Order, error) {
	var order model.Order

	err := or.db.First(order, id).Error 
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("order not found")
	} else if err != nil {
		return nil, err
	}

	err = or.db.Model(&order).Updates(updateOrder).Error
	if err != nil {
		return nil, err
	}

	err = or.db.Preload("User").Preload("OrderItems").First(&order, id).Error
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (or *orderRepository) DeleteOrder(id int) error {
	err := or.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("order_id = ?", id).Delete(&model.OrderItem{}).Error; err != nil {
			return err
		}

		if err := tx.Delete(&model.Order{}, id).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func NewOrderRepository(db *gorm.DB) *orderRepository {
	return &orderRepository{db: db}
}