package activities

import (
	"context"
	"errors"
	"ordermanagement/inventory"
	"ordermanagement/utils"

	"go.temporal.io/sdk/activity"
)

/* Ship
 * This activity sends an update to the warehouse to ship the order
 *
 * Takes a context.Context and an orderID and paymentInfo as parameters.
 * Returns a shipment confirmation and error.
 */
func Ship(ctx context.Context, order inventory.Order) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Ship Activity started")

	// simulate a random error
	if utils.IsError() {
		return "", errors.New("RANDOM CONFIRMING SHIPMENT ERROR: FELL OUT OF THE TRUCK")
	}

	return "CONFIRMED", nil
}
