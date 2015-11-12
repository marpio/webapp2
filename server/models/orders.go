package models

import (
	"fmt"
	. "github.com/ahmetalpbalkan/go-linq"
	"github.com/marpio/webapp2/server/constants/orderstate"
	vm "github.com/marpio/webapp2/server/viewmodels"
	"github.com/satori/go.uuid"
	"time"
)

type OrderRow struct {
	ID        string    `db:"id"`
	UserID    string    `db:"userId"`
	State     int       `db:"state"`
	CreatedAt time.Time `db:"createdAt"`
	ChangedAt time.Time `db:"changedAt"`
}

type ProductOrderRow struct {
	ID        string `db:"id"`
	ProductID string `db:"productId"`
	Amount    int    `db:"amount"`
}

type ServiceOrderRow struct {
	ID                string    `db:"id"`
	ServiceProviderID string    `db:"serviceProviderId"`
	ServiceDate       time.Time `db:"serviceDate"`
}

type orderTable struct {
	Base
	sot serviceOrderTable
	pot productOrderTable
}
type serviceOrderTable struct {
	Base
}
type productOrderTable struct {
	Base
}

func newOrderTable(ds *datastore) orderTable {
	o := orderTable{
		Base: Base{
			db:    ds.db,
			table: "orders",
		},
		sot: serviceOrderTable{
			Base: Base{
				db:    ds.db,
				table: "serviceOrders",
			},
		},
		pot: productOrderTable{
			Base: Base{
				db:    ds.db,
				table: "productOrders",
			},
		},
	}

	return o
}

func (ds *datastore) CreateOrder(order *vm.OrderDto) (*vm.OrderDto, error) {
	o := newOrderTable(ds)
	tx, _ := ds.db.Beginx()
	dataOrder := make(map[string]interface{})
	dataOrder["id"] = uuid.NewV4().String()
	dataOrder["userId"] = order.UserID
	dataOrder["state"] = orderstate.NEW
	dataOrder["createdAt"] = time.Now()
	dataOrder["changedAt"] = time.Now()
	_, err := o.InsertIntoTable(tx, dataOrder)
	if err != nil {
		return nil, err
	}

	dataServiceOrder := make(map[string]interface{})
	dataServiceOrder["id"] = uuid.NewV4().String()
	dataServiceOrder["orderId"] = dataOrder["id"]
	dataServiceOrder["serviceProviderId"] = order.ServiceOrder.ServiceProviderID
	dataServiceOrder["serviceDate"] = order.ServiceOrder.ServiceDate

	_, err = o.sot.InsertIntoTable(tx, dataServiceOrder)
	if err != nil {
		return nil, err
	}
	for _, v := range order.ProductOrders {
		dataProductOrder := make(map[string]interface{})
		dataProductOrder["orderId"] = dataOrder["id"]
		dataProductOrder["id"] = uuid.NewV4().String()
		dataProductOrder["productId"] = v.ProductID
		dataProductOrder["amount"] = v.Amount

		_, err = o.pot.InsertIntoTable(tx, dataProductOrder)
		if err != nil {
			return nil, err
		}
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (ds *datastore) GetOrderByID(id string) (*vm.OrderDto, error) {
	o := newOrderTable(ds)
	order := &OrderRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE id=$1", o.table)
	err := ds.db.Get(order, query, id)
	if err != nil {
		return nil, err
	}
	serviceOrder := &ServiceOrderRow{}
	querySo := fmt.Sprintf("SELECT * FROM %v WHERE orderId=$1", o.sot.table)
	err = ds.db.Get(order, querySo, id)
	if err != nil {
		return nil, err
	}
	productOrders := &[]ProductOrderRow{}
	queryPo := fmt.Sprintf("SELECT * FROM %v WHERE orderId=$1", o.pot.table)
	err = ds.db.Select(productOrders, queryPo, id)
	if err != nil {
		return nil, err
	}

	pos, err := From(*productOrders).Select(func(po T) (T, error) {
		return vm.ProductOrderDto{
			ID:        po.(ProductOrderRow).ID,
			ProductID: po.(ProductOrderRow).ProductID,
			Amount:    po.(ProductOrderRow).Amount,
		}, nil
	}).Results()
	if err != nil {
		return nil, err
	}

	productOrderDtos := make([]vm.ProductOrderDto, len(pos))
	for i, p := range pos {
		productOrderDtos[i] = p.(vm.ProductOrderDto)
	}

	res := &vm.OrderDto{
		ID:            order.ID,
		UserID:        order.UserID,
		State:         order.State,
		ProductOrders: productOrderDtos,
	}
	if serviceOrder != nil {
		res.ServiceOrder = vm.ServiceOrderDto{
			ID:                serviceOrder.ID,
			ServiceProviderID: serviceOrder.ServiceProviderID,
			ServiceDate:       serviceOrder.ServiceDate,
		}
	}
	return res, nil
}
