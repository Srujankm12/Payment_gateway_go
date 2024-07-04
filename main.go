package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/razorpay/razorpay-go"
)

// init is used to load environment variables from the .env file
func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
}

// homePage renders the home page with the payment button
func homePage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("app.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
	}
}

// createOrder creates a new order using Razorpay API and returns the order details as JSON
func createOrder(w http.ResponseWriter, r *http.Request) {
	// Initialize Razorpay client with key ID and secret from environment variables
	client := razorpay.NewClient(os.Getenv("RAZORPAY_KEY_ID"), os.Getenv("RAZORPAY_KEY_SECRET"))

	// Define order data
	data := map[string]interface{}{
		"amount":   1000,
		"currency": "INR",
		"receipt":  "receipt#1",
	}

	// Create the order
	order, err := client.Order.Create(data, nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating order: %v", err), http.StatusInternalServerError)
		return
	}

	// Convert order to JSON format
	orderJSON, err := json.Marshal(order)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error marshalling order: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(orderJSON)
}

func main() {
	// Define the route for the home page
	http.HandleFunc("/", homePage)
	// Define the route for creating an order
	http.HandleFunc("/createOrder", createOrder)

	// Start the server on port 8080
	port := ":8080"
	fmt.Printf("Server starting on port %s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
