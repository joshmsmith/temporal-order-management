package workflows

import (
	"ordermanagement/activities"
	"ordermanagement/inventory"
	"time"

	"go.temporal.io/sdk/workflow"
)

// InventoryWorkflow is a function that handles the inventory workflow.
// It takes a workflow context as input and returns an error if any.
func InventoryWorkflow(ctx workflow.Context, order inventory.Order) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("InventoryWorkflow")

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	// 	-- process_order(order):
	//-- check_fraud(order.order_id, order.payment_info)
	//-- prepare_shipment(order)
	//-- charge_confirm = charge(order.order_id, order.payment_info) // worker dies here
	//-- shipment_confirmation = ship(order)

	err := workflow.ExecuteActivity(ctx, activities.UpdateInventoryActivity, order).Get(ctx, nil)
	if err != nil {
		logger.Error("Activity failed.", "Error", err)
		return err
	}

	err = workflow.ExecuteActivity(ctx, activities.SupplierOrderActivity, order.Item, 1000).Get(ctx, nil)

	logger.Info("InventoryWorkflow completed.")

	return nil

}
