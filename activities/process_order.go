package activities

import (
	"context"
	"errors"
	"ordermanagement/inventory"
	"ordermanagement/utils"

	"go.temporal.io/sdk/activity"
)

/* ProcessOrder
 *   This activity processes an order, and demonstrates Idempotence: if an order already exists, do nothing, so it can be called multiple times with no effect
 *   1. checks if the order exists
 *   2. validates that there's enough inventory for the order
 *   2. updates the inventory (stock) based on the given order and writes the order as existing
 *
 * Takes a context.Context and an inventory.Order as parameters.
 * Returns an indicator if there is a duplicate oder and an error.
 */
func ProcessOrder(ctx context.Context, order inventory.Order) (bool, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("ProcessOrder started")

	// Idempotence: If order already processed, do nothing
	orderExists := inventory.SearchOrder(order.OrderID)
	if orderExists {
		logger.Info("ProcessOrder: Order", order.OrderID, "Exists!")
		return true, nil // nothing to do
	}

	inStock, err := inventory.GetInStock(order.Item)
	if err != nil {
		logger.Error("ProcessOrder: checking stock failed", "Error", err)
		return false, err
	}

	if inStock < order.Quantity {
		logger.Error("ProcessOrder: not enough stock", "Error", err)
		return false, errors.New("NOT ENOUGH STOCK")
	}

	// simulate a random error
	if utils.IsError() {
		logger.Error("ProcessOrder: checking stock failed for random reason", "Error", err)
		return false, errors.New("ERROR CHECKING STOCK")
	}

	err = inventory.UpdateStock(order.OrderID, order.Item, inStock-order.Quantity)
	if err != nil {
		return false, err
	}

	return false, nil
}
