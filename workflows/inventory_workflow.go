package workflows

import (
	"ordermanagement/activities"
	"ordermanagement/inventory"
	"time"

	"go.temporal.io/sdk/workflow"

	u "ordermanagement/utils"
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
	err := workflow.ExecuteActivity(ctx, activities.ProcessOrder, order).Get(ctx, nil)
	if err != nil {
		logger.Error("ProcessOrder activity failed.", "Error", err)
		return err
	}
	//-- check_fraud(order.order_id, order.payment_info)
	err = workflow.ExecuteActivity(ctx, activities.CheckFraud, order.OrderID, order.PaymentInfo).Get(ctx, nil)
	if err != nil {
		logger.Error("CheckFraud activity failed.", "Error", err)
		return err
	}
	//-- prepare_shipment(order)
	err = workflow.ExecuteActivity(ctx, activities.PrepareShipment, order).Get(ctx, nil)
	if err != nil {
		logger.Error("PrepareShipment activity failed.", "Error", err)
		return err
	}
	//-- charge_confirm = charge(order.order_id, order.payment_info) // worker dies here
	chargeConfirm := "UNCONFIRMED"
	err = workflow.ExecuteActivity(ctx, activities.ChargeConfirm, order.OrderID, order.PaymentInfo).Get(ctx, &chargeConfirm)
	if err != nil {
		logger.Error("ChargeConfirm activity failed.", "Error", err)
		return err
	}

	// For demo - sleep between activities so you can kill the worker
	delay := 15
	logger.Debug("InventoryWorkflow: (DEBUG) Sleeping between activity calls -")
	logger.Info(u.ColorGreen, "InventoryWorkflow:", u.ColorBlue, "workflow.Sleep duration", delay, "seconds", u.ColorReset)
	workflow.Sleep(ctx, time.Duration(delay)*time.Second)

	//-- shipment_confirmation = ship(order)
	shipmentConfirmation := "UNCONFIRMED"
	err = workflow.ExecuteActivity(ctx, activities.Ship, order).Get(ctx, &shipmentConfirmation)
	if err != nil {
		logger.Error("ChargeConfirm activity failed.", "Error", err)
		return err
	}

	//bonus: order more stuff if needed
	err = workflow.ExecuteActivity(ctx, activities.SupplierOrderActivity, order.Item, 10000).Get(ctx, nil)

	logger.Info("InventoryWorkflow completed. Charge Status:", chargeConfirm, ", Shipment Status: ", shipmentConfirmation, ".")

	return nil

}
