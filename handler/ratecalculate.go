package handlers // Make sure this matches your folder structure

import (
	"database/sql"
	db "dod-backend/Database"
	"encoding/json"
	"net/http"
	"strconv"
	// "dod-backend/Database" // Assuming you have your DB imported here
)

// 1. Capital 'G' makes it available to the routes package
func GetRateHandler(w http.ResponseWriter, r *http.Request) {
	// 2. Standard way to get a query parameter
	distanceParam := r.URL.Query().Get("distance")
	actualDistance, err := strconv.ParseFloat(distanceParam, 64)

	w.Header().Set("Content-Type", "application/json") // Set response type to JSON

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Invalid distance provided"}`))
		return
	}

	query := `SELECT ratePerKM FROM distanceRate WHERE ? >= min_distance AND ? <= max_distance LIMIT 1`
	var ratePerKM int

	// Replace 'db.DB' with whatever your database connection variable is
	err = db.DB.QueryRow(query, actualDistance, actualDistance).Scan(&ratePerKM)

	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"error": "Rate not defined for this distance"}`))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "Database query failed"}`))
		}
		return
	}

	// 3. Send the successful JSON response
	json.NewEncoder(w).Encode(map[string]interface{}{
		"distance":    actualDistance,
		"ratePerKM":   ratePerKM,
		"total price": float64(ratePerKM) * actualDistance,
	})
}
