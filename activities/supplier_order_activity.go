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
	//TODO make these random errors more better fake errors
	if utils.IsError() {
		return errors.New("RANDOM ERROR")
	}
	if inStock < 10 {
		// Call supplier API and update inventory
		if utils.IsError() {
			//TODO make these random errors more better fake errors
			return errors.New("RANDOM ERROR")
		}
		inventory.SupplierOrder(quantity)
		return nil
	}
	return nil
}
