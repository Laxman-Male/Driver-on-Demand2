package handlers

import (
	db "dod-backend/Database"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// GetPaymentDetails fetches the amount for a booking before paying
func GetPaymentDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	bookingID := r.URL.Query().Get("id")
	if bookingID == "" {
		http.Error(w, "Booking ID missing", http.StatusBadRequest)
		return
	}

	var fare float64
	query := `SELECT total_fare FROM RideBookingstbl WHERE BookingID = ?`
	err := db.DB.QueryRow(query, bookingID).Scan(&fare)

	if err != nil {
		http.Error(w, "Booking not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"bookingId": bookingID,
		"amount":    fare,
		"upiId":     "yourname@okaxis", // Replace with your actual UPI for testing
	})
}

// ConfirmPayment handles the mock 'Success' from frontend
func ConfirmPayment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var p struct {
		BookingID int     `json:"bookingId"`
		Amount    float64 `json:"amount"`
		Method    string  `json:"method"`
	}

	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// 1. Generate a Mock Transaction ID
	txID := fmt.Sprintf("TXN-%d", time.Now().Unix())

	// 2. Start a Database Transaction to ensure both tables update or none do
	tx, err := db.DB.Begin()
	if err != nil {
		http.Error(w, "Transaction error", http.StatusInternalServerError)
		return
	}

	// 3. Insert into Payments table
	_, err = tx.Exec("INSERT INTO Paymentstbl (BookingID, Amount, PaymentMethod, TransactionID, PaymentStatus) VALUES (?, ?, ?, ?, 'Success')",
		p.BookingID, p.Amount, p.Method, txID)

	if err != nil {
		tx.Rollback()
		http.Error(w, "Payment record failed", http.StatusInternalServerError)
		return
	}

	// 4. Update Ride Status to 'Completed'
	_, err = tx.Exec("UPDATE RideBookingstbl SET status = 'Completed', updated_at = NOW() WHERE BookingID = ?", p.BookingID)

	if err != nil {
		tx.Rollback()
		http.Error(w, "Ride status update failed", http.StatusInternalServerError)
		return
	}

	// 5. Commit the changes
	tx.Commit()

	json.NewEncoder(w).Encode(map[string]string{
		"message":       "Payment successful and ride completed!",
		"transactionId": txID,
	})
}
