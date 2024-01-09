package main

import (
	"errors"
	"fmt"
	"ordermanagement/inventory"
	"strconv"
	"time"

	"log"
	"os"

	"math/rand"

	u "ordermanagement/utils"
)

/* main
 * entry point - set up and start the  process with a default order setup and random order number
 * optional command line arguments:
 *   order number: int
 *   product number: int (only 123456 is supported)
 *   quantity to order: int
 *   payment info: string
 */
func main() {

	// Set order values
	order := inventory.Order{
		OrderID:     fmt.Sprintf("order-%d", rand.Intn(99999)),
		Item:        "123456",
		Quantity:    999,
		PaymentInfo: "VISA-12345",
	}

	// collect command line args if valid
	args := os.Args[1:]
	if len(args) >= 4 {
		if i, err := strconv.Atoi(args[2]); err == nil {
			order = inventory.Order{
				OrderID:     fmt.Sprintf("order-%s", args[0]),
				Item:        args[1],
				Quantity:    i,
				PaymentInfo: args[3],
			}
		}
	}

	log.Print(order.OrderID, " - called, Order Details: ", order)
	// Checking actual values
	inStock, err := inventory.GetInStock(order.Item)
	log.Println(order.OrderID, "- Product #", order.Item, "stock:", inStock)

	result, err := ProcessOrder(order)
	

	if err != nil {
		log.Fatalf("%s - %sProcessing Order returned failure:%s %v", order.OrderID, u.ColorRed, u.ColorReset, err)
	} else {
		log.Printf("%s - %sProcessing Order completed:%s Result: %s", order.OrderID, u.ColorGreen, u.ColorReset, result)
	}

	// Get product stock
	inStock, err = inventory.GetInStock(order.Item)
	log.Printf("%s - Current Product %s stock: %d \n", order.OrderID, order.Item, inStock)

}

// ProcessOrder is a function that handles the inventory workflow.
// It takes a workflow context as input and returns an error if any.
func ProcessOrder(order inventory.Order) (string, error) {
	
	log.Println(order.OrderID, "- Starting ProcessOrder, Order Details: ", order)


	//-- check_fraud(order.order_id, order.payment_info)
	err := CheckFraud(order.OrderID, order.PaymentInfo)
	
	if err != nil {
		log.Println("CheckFraud process failed.", "Error", err)
		return "", err
	}

	// For demo - sleep between activities so you can kill the worker
	delay := 0
	if order.OrderID == "order-37005" {
		delay = 15
		log.Printf("ProcessOrder: Sleeping between process calls -")
		log.Println(u.ColorGreen, "ProcessOrder:", u.ColorBlue, "workflow.Sleep duration", delay, "seconds", u.ColorReset)
		time.Sleep(time.Duration(delay)*time.Second)
	}	

	//-- prepare_shipment(order)
	err = PrepareShipment( order)
	if err != nil {
		log.Println("PrepareShipment method failed.", "Error", err)
		return "", err
	}

	//-- charge_confirm = charge(order.order_id, order.payment_info) // worker dies here
	charge := "UNCONFIRMED"
	charge, err = Charge(order.OrderID, order.PaymentInfo)
	if err != nil {
		log.Println("Charge process failed.", "Error", err)
		return "", err
	}

	//-- shipment_confirmation = ship(order)
	shipmentConfirmation := "UNCONFIRMED"
	shipmentConfirmation, err = Ship(order)
	if err != nil {
		log.Println("Shipment process failed.", "Error", err)
		return "", err
	}

	log.Println("ProcessOrder completed. Charge Status:", charge, ", Shipment Status:", shipmentConfirmation, ".")

	return "Order Managed", nil
}

/* PrepareShipment
 *   This process prepares an order, and demonstrates Idempotence: if an order already exists, do nothing, so it can be called multiple times with no effect
 *   1. checks if the order exists
 *   2. validates that there's enough inventory for the order
 *   2. updates the inventory (stock) based on the given order and writes the order as existing
 *
 * Takes an inventory.Order as parameters.
 * Returns  an error.
 */
 func PrepareShipment(order inventory.Order) (error) {
	
	log.Printf("PrepareShipment Process started")

	// Idempotence: If order already processed, do nothing
	//orderExists := inventory.SearchOrder(order.OrderID)
	//if orderExists {
	//	log.Info("Order Exists!")
	//	return true, nil // nothing to do
	//}

	inStock, err := inventory.GetInStock(order.Item)
	if err != nil {
		return err
	}

	if inStock < order.Quantity {
		return errors.New("NOT ENOUGH STOCK")
	}

	// simulate a random error
	if u.IsError() {
		return errors.New("RANDOM ERROR PREPARING SHIPMENT: WAREHOUSE CAUGHT ON FIRE AND ALL THE BOXES BURNT UP")
	}

	err = inventory.UpdateStock(order.OrderID, order.Item, inStock-order.Quantity)
	if err != nil {
		return err
	}

	return  nil
}

/* Charge
 * This process confirms payment
 *
 * Takes an orderID and paymentInfo as parameters.
 * Returns an charge confirmation and error.
 */
 func Charge(orderID string, paymentInfo string) (string, error) {
	log.Printf("Charge process started")

	// pretend to charge, sometimes error
	if u.IsError() {
		return "", errors.New("RANDOM CONFIRMING CHARGE ERROR: NO MONEY")
	}

	return "CONFIRMED", nil
}


/* CheckFraud
 * This process validates that the payment info isn't fraudulent in a mocked out way
 *
 * Takes a an orderID and paymentInfo as parameters.
 * Returns an error.
 */
 func CheckFraud(orderID string, paymentInfo string) error {
	
	log.Printf("CheckFraud process started")

	// pretend to check for fraud, sometimes error
	if u.IsError() {
		return errors.New("RANDOM ERROR CHECKING FRAUD: THESE GUYS ARE PIRATES")
	}

	return nil
}


/* Ship
 * This process sends an update to the warehouse to ship the order
 *
 * Takes an orderID and paymentInfo as parameter.
 * Returns a shipment confirmation and error.
 */
 func Ship(order inventory.Order) (string, error) {
	log.Printf("Ship process started")

	// pretend to request the warehouse ship the order
	if u.IsError() {
		return "", errors.New("RANDOM CONFIRMING SHIPMENT ERROR: WAREHOUSE AUTOMATION SYSTEM DOWN")
	}
	log.Println("Shipment request successful for", order.OrderID, "- shipping", order.Quantity, "of item #", order.Item, "to requestor!")

	return "CONFIRMED", nil
}



// SupplierOrder call supplier API to order new product.
//
// It takes an item string as parameter.
// It returns an error.
func SupplierOrder(item string, quantity int) error {
	log.Printf("SupplierOrder started")

	inStock, err := inventory.GetInStock(item)
	if err != nil {
		return err
	}

	if u.IsError() {
		return errors.New("RANDOM ERROR CHECKING STOCK FOR RE-ORDER")
	}
	log.Println("SupplierOrder: Stock Level for", item, "is at:", inStock)
	if inStock < 5000 {
		// Call supplier API and update inventory
		log.Println("SupplierOrder: Stock Level less than minimum required, re-ordering up to", quantity)
		if u.IsError() {

			return errors.New("RANDOM ERROR TELLING SUPPLIER TO SEND US MORE STUFF")
		}
		inventory.SupplierOrder(quantity)
		return nil
	}
	return nil
}
