package activities

import (
	"context"
	"errors"
	"ordermanagement/inventory"
	"ordermanagement/utils"

	"go.temporal.io/sdk/activity"
)

/* PrepareShipment
 *   This activity prepares an order, and demonstrates Idempotence: if an order already exists, do nothing, so it can be called multiple times with no effect
 *   1. checks if the order exists
 *   2. validates that there's enough inventory for the order
 *   2. updates the inventory (stock) based on the given order and writes the order as existing
 *
 * Takes a context.Context and an inventory.Order as parameters.
 * Returns an error.
 */
func PrepareShipment(ctx context.Context, order inventory.Order) error {
	logger := activity.GetLogger(ctx)
	logger.Info("PrepareShipment Activity started")

	// Idempotence: If order already processed, do nothing
	orderExists := inventory.SearchOrder(order.OrderID)
	if orderExists {
		logger.Info("Order Exists!")
		return nil // nothing to do
	}

	inStock, err := inventory.GetInStock(order.Item)
	if err != nil {
		return err
	}

	if inStock < order.Quantity {
		return errors.New("NOT ENOUGH STOCK")
	}

	// simulate a random error
	if utils.IsError() {
		return errors.New("RANDOM ERROR PREPARING SHIPMENT: WAREHOUSE CAUGHT ON FIRE AND ALL THE BOXES BURNT UP")
	}

	err = inventory.UpdateStock(order.OrderID, order.Item, inStock-order.Quantity)
	if err != nil {
		return err
	}

	return nil
}
