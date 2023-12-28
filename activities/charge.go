package activities

import (
	"context"
	"errors"
	"ordermanagement/utils"

	"go.temporal.io/sdk/activity"
)

/* Charge
 * This activity confirms payment
 *
 * Takes a context.Context and an orderID and paymentInfo as parameters.
 * Returns an charge confirmation and error.
 */
func Charge(ctx context.Context, orderID string, paymentInfo string) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Charge Activity started")

	// pretend to charge, sometimes error
	if utils.IsError() {
		return "", errors.New("RANDOM CONFIRMING CHARGE ERROR: NO MONEY")
	}

	return "CONFIRMED", nil
}
