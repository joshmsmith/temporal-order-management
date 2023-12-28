package inventory

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Inventory struct {
	OrderID   string `json:"orderID"`
	ProductID string `json:"productID"`
	InStock   int    `json:"inStock"`
}

// ReadJSON reads a JSON file and returns the unmarshalled data as an Inventory struct.
//
// It takes a filename string as a parameter.
// It returns an Inventory struct and an error.
func ReadJSON(filename string) (Inventory, error) {
	file, err := os.Open(filename)
	if err != nil {
		return Inventory{}, err
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return Inventory{}, err
	}

	var data Inventory
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return Inventory{}, err
	}

	return data, nil
}

// SearchOrder searches for an order in the database by its order ID.
//
// Parameters:
// - orderID: a string representing the order ID to search for.
//
// Returns:
// - bool: true if the order is found, false otherwise.
func SearchOrder(orderID string) bool {

	data, err := ReadJSON(os.Getenv("DATABASE"))
	if err != nil {
		fmt.Println("Error reading JSON:", err)
		return false
	}

	if data.OrderID == orderID {
		return true
	}
	return false
}

// GetInStock retrieves the quantity of a product that is currently in stock.
//
// It takes a string parameter, productID, which represents the unique identifier of the product.
// The function returns an integer value representing the quantity of the product that is currently in stock.
func GetInStock(productID string) (int, error) {

	data, err := ReadJSON(os.Getenv("DATABASE"))
	if err != nil {
		fmt.Println("Error reading JSON:", err)
		return 0, err
	}
	if data.ProductID == productID {
		return data.InStock, err
	}
	return 0, nil
}

// UpdateStock updates the stock of a product for a given order in the database.
//
// Parameters:
// - orderID: the ID of the order.
// - productID: the ID of the product.
// - inStock: the new stock value for the product.
//
// Return type: error.
func UpdateStock(orderID string, productID string, inStock int) error {

	data, err := ReadJSON(os.Getenv("DATABASE"))
	if err != nil {
		fmt.Println("Error reading JSON:", err)
		return err
	}
	data.OrderID = orderID
	data.InStock = inStock
	data.ProductID = productID
	UpdateJSON(os.Getenv("DATABASE"), data)
	return nil
}

// UpdateJSON updates the given JSON file with the provided data.
//
// The function returns an error if any error occurs during the file opening,
// encoding, or closing process. Otherwise, it returns nil.
func UpdateJSON(filename string, data Inventory) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ") // Pretty print the JSON
	err = encoder.Encode(data)
	if err != nil {
		return err
	}

	return nil
}

func SupplierOrder(quantity int) error {

	data, err := ReadJSON(os.Getenv("DATABASE"))
	if err != nil {
		fmt.Println("Error reading JSON:", err)
		return err
	}
	GetInStock(data.ProductID)
	// how much should we order?
	toOrder := quantity - data.InStock

	// order that much
	fmt.Println("Ordering from supplier: ", toOrder)
	data.InStock = data.InStock + toOrder
	UpdateJSON(os.Getenv("DATABASE"), data)
	return nil
}
