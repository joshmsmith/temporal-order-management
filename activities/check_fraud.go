package activities

import (
	"context"
	"errors"
	"ordermanagement/utils"

	"go.temporal.io/sdk/activity"
)

/* CheckFraud
 * This activity validates that the payment info isn't fraudulent in a mocked out way
 *
 * Takes a context.Context and an orderID and paymentInfo as parameters.
 * Returns an error.
 */
func CheckFraud(ctx context.Context, orderID string, paymentInfo string) error {
	logger := activity.GetLogger(ctx)
	logger.Info("CheckFraud Activity started")

	// simulate a random error
	if utils.IsError() {
		return errors.New("RANDOM ERROR CHECKING FRAUD ACTIVITY: THESE GUYS ARE PIRATES")
	}

	return nil
}
