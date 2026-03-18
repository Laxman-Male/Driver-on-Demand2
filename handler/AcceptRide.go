// handlers/rides.go
package handlers

import (
	db "dod-backend/Database"
	"fmt"
	"net/http"
)

// AcceptRideHandler updates status to 'InProcess'
func AcceptRideHandler(w http.ResponseWriter, r *http.Request) {
	// Enable CORS for this specific handler if not using middleware
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	// Get BookingID from query parameter: /accept-ride?id=123
	bookingID := r.URL.Query().Get("id")
	if bookingID == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Booking ID is required"}`))
		return
	}
	fmt.Println("===", bookingID)
	// Update query
	query := `UPDATE RideBookingstbl SET status = 'InProcess', updated_at = NOW() WHERE BookingID = ?`
	_, err := db.DB.Exec(query, bookingID)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Failed to update database"}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Ride status updated to InProcess"}`))
}
