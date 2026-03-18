package main

import (
	db "dod-backend/Database"
	handlers "dod-backend/handler"
	"fmt"
	"log"
	"net/http"

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

	// CORS middleware
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:4200"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})

	fmt.Println("Server running on :8000")
	http.ListenAndServe(":8000", c.Handler(mux))
}
