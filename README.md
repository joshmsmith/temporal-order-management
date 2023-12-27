# idempotence

This project demonstrates how to manage idempotence in a Temporal.io project. Check out the full blog post here: [Idempotence in Temporal.io, a Look into Technical Architectures](https://medium.com/@ps.augereau/idempotence-in-temporal-io-a-look-into-technical-architectures-11d20a0fc860)

When a customer orders a product, the inventory is decrease and if the stock is below 2 and order to the supplier is made and the stock is incremented.
 
This project utilises a workflow named **"Inventory Workflow"**, and this workflow invokes two activities:

 1. UpdateInventoryActivity. The inventory is updated with the decreasing quantity.
 2. Supplier Order. If the stock for the product drops below 2, an order is placed with the supplier and the stock is updated.

**Idempotence** ensures that when a workflow or an activity is retried by Temporal (due to an error for example), the system remains unchanged. This means:
 - Updating the inventory with the same orderID results in no change.
 - If an order has already been sent to the supplier, no additional orders are sent.

**All this rules have to be managed by the developer.**

The activities introduce random errors (this is for testing Temporal's retry capabilities and idempotency)
The code in `starter.go`` demonstrates the workflow. Initially, product 123456 has 1000 units. A customer orders 999 units, and the supplier credits 1000 units. Therefore, the final stock for product 123456 is 1001 units.

At the beginning, the database looks like this:

```json
{
    "orderID": "000",
    "productID": "123456",
    "inStock": 1000
}
```

At the end the database looks like this:

```json
{
    "orderID": "A12",
    "productID": "123456",
    "inStock": 1001
}
```


# Setup
You need to have a Temporal server up and running.

CLone the repo :
```shell
git clone git@github.com:PierreSylvain/idempotence.git
```

Copy the .env-dist file into .env and change the values as needed :

```shell
TEMPORAL_URL=localhost:7233
TEMPORAL_NAMESPACE=default
TASKQUEUE=inventory-task-queue
DATABASE=database/inventory.json
```

The `DATABASE` parameter in the name where to store the data.

Then start the worker :
```shell
go run worker/main.go
```

And finally, test with :
```shell 
go run starter/main.go
```



