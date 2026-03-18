// package main

// import (
// 	db "dod-backend/Database"
// 	handlers "dod-backend/handler"
// 	"fmt"
// 	"log"
// 	"net/http"

// 	"github.com/rs/cors"
// )

// func main() {
// 	// Initialize DB once
// 	if err := db.Init(); err != nil {
// 		log.Fatal("Failed to connect to database:", err)
// 	}

// 	mux := http.NewServeMux()
// 	mux.HandleFunc("/login", handlers.LoginUser)
// 	mux.HandleFunc("/get-rate", handlers.GetRateHandler)
// 	mux.HandleFunc("/book-ride/", handlers.BookRideHandler)
// 	mux.HandleFunc("/get-pending-rides", handlers.GetPendingRidesForDriver)
// 	mux.HandleFunc("/accept-ride", handlers.AcceptRideHandler)
// 	mux.HandleFunc("/get-completed-rides", handlers.GetCompletedRidesByMobile)
// 	mux.HandleFunc("/complete-ride", handlers.CompleteRideHandler)
// 	mux.HandleFunc("/get-all-bookings", handlers.GetAllUserBookings)

// 	//payment gateway
// 	mux.HandleFunc("/get-payment-info", handlers.GetPaymentDetails)
// 	mux.HandleFunc("/process-payment", handlers.ConfirmPayment)

// 	mux.HandleFunc("/driver-completed-history", handlers.GetDriverCompletedRides)

// 	// CORS middleware
// 	c := cors.New(cors.Options{
// 		AllowedOrigins:   []string{"http://localhost:4200"},
// 		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
// 		AllowedHeaders:   []string{"Content-Type"},
// 		AllowCredentials: true,
// 	})

// 	fmt.Println("Server running on :8000")
// 	http.ListenAndServe(":8000", c.Handler(mux))
// }

package main

import (
	db "dod-backend/Database"
	handlers "dod-backend/handler"
	"fmt"
	"log"
	"net/http"
	"os" // Added this to read cloud variables

	"github.com/rs/cors"
)

func main() {
	// Initialize DB once
	if err := db.Init(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/login", handlers.LoginUser)
	mux.HandleFunc("/get-rate", handlers.GetRateHandler)
	mux.HandleFunc("/book-ride/", handlers.BookRideHandler)
	mux.HandleFunc("/get-pending-rides", handlers.GetPendingRidesForDriver)
	mux.HandleFunc("/accept-ride", handlers.AcceptRideHandler)
	mux.HandleFunc("/get-completed-rides", handlers.GetCompletedRidesByMobile)
	mux.HandleFunc("/complete-ride", handlers.CompleteRideHandler)
	mux.HandleFunc("/get-all-bookings", handlers.GetAllUserBookings)

	//payment gateway
	mux.HandleFunc("/get-payment-info", handlers.GetPaymentDetails)
	mux.HandleFunc("/process-payment", handlers.ConfirmPayment)

	mux.HandleFunc("/driver-completed-history", handlers.GetDriverCompletedRides)

	// 1. FIX CORS FOR CLOUD
	// Read the live frontend URL from Render, fallback to localhost for testing
	allowedOrigin := os.Getenv("FRONTEND_URL")
	if allowedOrigin == "" {
		allowedOrigin = "http://localhost:4200"
	}

	// CORS middleware
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{allowedOrigin, "http://localhost:4200"}, // Allows both live and local
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"}, // Added Authorization just in case
		AllowCredentials: true,
	})

	// 2. FIX PORT FOR CLOUD
	// Read the dynamic port from Render, fallback to 8000 for local testing
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	fmt.Println("Server running on port :" + port)

	// Start the server using the dynamic port
	err := http.ListenAndServe(":"+port, c.Handler(mux))
	if err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
