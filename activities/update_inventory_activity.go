package activities

import (
	"context"
	"errors"
	"ordermanagement/inventory"
	"ordermanagement/utils"

	"go.temporal.io/sdk/activity"
)

// UpdateInventoryActivity updates the inventory based on the given order.
//
// Takes a context.Context and an inventory.Order as parameters.
// Returns an error.
func UpdateInventoryActivity(ctx context.Context, order inventory.Order) error {
	logger := activity.GetLogger(ctx)
	logger.Info("UpdateInventoryActivity started")

	// Idempotence: If order already processed, do nothing
	orderExists := inventory.SearchOrder(order.OrderID)
	if orderExists {
		return nil
	}

	inStock, err := inventory.GetInStock(order.Item)
	if err != nil {
		return err
	}

	//TODO make these random errors more better fake errors
	if utils.IsError() {
		return errors.New("RANDOM ERROR")
	}
	err = inventory.UpdateStock(order.OrderID, order.Item, inStock-order.Quantity)
	if err != nil {
		return err
	}
	if utils.IsError() {
		//TODO make these random errors more better fake errors
		return errors.New("RANDOM ERROR")
	}
	return nil
}
