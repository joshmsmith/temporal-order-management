package activities

import (
	"context"
	"errors"
	"ordermanagement/inventory"
	"ordermanagement/utils"

	"go.temporal.io/sdk/activity"
)

// SupplierOrderActivity call supplier API to order new product.
//
// It takes a context.Context and an item string as parameters.
// It returns an error.
func SupplierOrderActivity(ctx context.Context, item string, quantity int) error {
	logger := activity.GetLogger(ctx)
	logger.Info("SupplierOrderActivity started")

	inStock, err := inventory.GetInStock(item)
	if err != nil {
		return err
	}

	if utils.IsError() {
		return errors.New("RANDOM ERROR CHECKING STOCK FOR RE-ORDER")
	}
	logger.Info("SupplierOrderActivity: Stock Level for", item, "is at:", inStock)
	if inStock < 5000 {
		// Call supplier API and update inventory
		logger.Info("SupplierOrderActivity: Stock Level less than minimum required, re-ordering up to", quantity)
		if utils.IsError() {

			return errors.New("RANDOM ERROR TELLING SUPPLIER TO SEND US MORE STUFF")
		}
		inventory.SupplierOrder(quantity)
		return nil
	}
	return nil
}
