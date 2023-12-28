package workflows

import (
	"ordermanagement/activities"
	"ordermanagement/inventory"
	"time"

	"go.temporal.io/sdk/workflow"

	u "ordermanagement/utils"
)

// ProcessOrder is a function that handles the inventory workflow.
// It takes a workflow context as input and returns an error if any.
func ProcessOrderAndValidateStock(ctx workflow.Context, order inventory.Order) (string, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info(order.OrderID, " - Starting ProcessOrderAndValidateStock, Order Details: ", order)

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	//-- check_fraud(order.order_id, order.payment_info)
	err := workflow.ExecuteActivity(ctx, activities.CheckFraud, order.OrderID, order.PaymentInfo).Get(ctx, nil)
	if err != nil {
		logger.Error("CheckFraud activity failed.", "Error", err)
		return "", err
	}

	// For demo - sleep between activities so you can kill the worker
	delay := 0
	if order.OrderID == "order-37005" {
		delay = 15
	}
	logger.Debug("ProcessOrderAndValidateStock: Sleeping between activity calls -")
	logger.Info(u.ColorGreen, "ProcessOrderAndValidateStock:", u.ColorBlue, "workflow.Sleep duration", delay, "seconds", u.ColorReset)
	workflow.Sleep(ctx, time.Duration(delay)*time.Second)

	// 	-- process_order(order):
	duplicate := false

	//-- prepare_shipment(order)
	err = workflow.ExecuteActivity(ctx, activities.PrepareShipment, order).Get(ctx, &duplicate)
	if err != nil {
		logger.Error("PrepareShipment activity failed.", "Error", err)
		return "", err
	}

	if duplicate {
		return "Duplicate order", nil
	}

	//-- charge_confirm = charge(order.order_id, order.payment_info) // worker dies here
	Charge := "UNCONFIRMED"
	err = workflow.ExecuteActivity(ctx, activities.Charge, order.OrderID, order.PaymentInfo).Get(ctx, &Charge)
	if err != nil {
		logger.Error("Charge activity failed.", "Error", err)
		return "", err
	}

	//-- shipment_confirmation = ship(order)
	shipmentConfirmation := "UNCONFIRMED"
	err = workflow.ExecuteActivity(ctx, activities.Ship, order).Get(ctx, &shipmentConfirmation)
	if err != nil {
		logger.Error("Charge activity failed.", "Error", err)
		return "", err
	}

	//bonus: order more stuff if needed
	err = workflow.ExecuteActivity(ctx, activities.SupplierOrderActivity, order.Item, 10000).Get(ctx, nil)

	logger.Info("ProcessOrderAndValidateStock completed. Charge Status:", Charge, ", Shipment Status: ", shipmentConfirmation, ".")

	return "Order Managed", nil

}
