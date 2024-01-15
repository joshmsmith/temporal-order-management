package main

import (
	"context"
	"fmt"
	"ordermanagement/inventory"
	"ordermanagement/workflows"
	"strconv"

	"log"
	"os"

	//"github.com/google/uuid"
	"math/rand"
	//"github.com/joho/godotenv"
	"go.temporal.io/sdk/client"

	u "ordermanagement/utils"
)

var OrderManagementTransferTaskQueueName = os.Getenv("ORDER_MANAGEMENT_TASK_QUEUE")

/* main
 * entry point - set up and start the workflow process with a default order setup and random order number
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
		Quantity:    13,
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

	// Load the Temporal Cloud from env
	clientOptions, err := u.LoadClientOptions(u.NoSDKMetrics)
	if err != nil {
		log.Fatalf(order.OrderID, "- Failed to load Temporal Cloud environment: %v", err)
	}
	log.Print(order.OrderID, " - connecting to Temporal server.")
	temporalClient, err := client.Dial(clientOptions)
	if err != nil {
		log.Fatalln(order.OrderID, "- Unable to create Temporal client: %v", err)
	}
	defer temporalClient.Close()

	// Workflow options
	workflowID := fmt.Sprintf("ordermanagement-wkfl-%s", order.OrderID)

	workflowOptions := client.StartWorkflowOptions{
		ID:        workflowID,
		TaskQueue: OrderManagementTransferTaskQueueName,
	}

	// Checking actual values
	inStock, err := inventory.GetInStock(order.Item)
	log.Println(order.OrderID, "- Product #", order.Item, "stock:", inStock)

	// Execute workflow
	log.Println(order.OrderID, "- Starting Order Management System workflow on", OrderManagementTransferTaskQueueName, "task queue")
	workflowExec, err := temporalClient.ExecuteWorkflow(context.Background(), workflowOptions, workflows.ProcessOrder, order)
	if err != nil {

		log.Fatalln(order.OrderID, "- %sError, Unable to execute workflow %v:%s", u.ColorRed, err, u.ColorReset)
	}
	log.Printf("%s - %sWorkflow started:%s (WorkflowID: %s, RunID: %s)", order.OrderID, u.ColorYellow, u.ColorReset, workflowExec.GetID(), workflowExec.GetRunID())

	// Wait for the workflow completion.
	var result string
	errWF := workflowExec.Get(context.Background(), &result)

	if errWF != nil {
		log.Fatalln(order.OrderID, "- %sWorkflow returned failure:%s %v", u.ColorRed, u.ColorReset, errWF)
	} else {
		log.Printf("%s - %sWorkflow completed:%s WorkflowID: %s, RunID: %s, Result: %s", order.OrderID, u.ColorGreen, u.ColorReset, workflowExec.GetID(), workflowExec.GetRunID(), result)
	}

	// Get product stock
	inStock, err = inventory.GetInStock(order.Item)
	log.Printf("%s - Current Product %s stock: %d \n", order.OrderID, order.Item, inStock)

}
