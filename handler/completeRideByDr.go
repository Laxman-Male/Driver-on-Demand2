// handlers/rides.go
package handlers

import (
	db "dod-backend/Database"
	"net/http"
)

// CompleteRideHandler updates status to 'Completed'
func CompleteRideHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	// Get BookingID from query parameter: /complete-ride?id=123
	bookingID := r.URL.Query().Get("id")
	if bookingID == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Booking ID is required"}`))
		return
	}

	// Update query to set status to 'Completed'
	query := `UPDATE RideBookingstbl SET status = 'Completed', updated_at = NOW() WHERE BookingID = ?`
	_, err := db.DB.Exec(query, bookingID)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Failed to complete ride in database"}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Ride marked as Completed"}`))
}
