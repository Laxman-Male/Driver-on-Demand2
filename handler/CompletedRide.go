package handlers

import (
	db "dod-backend/Database"
	"encoding/json"
	"net/http"
)

// GetCompletedRidesByMobile fetches history for a specific phone number
func GetCompletedRidesByMobile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// 1. Get mobile from query param: /get-completed-rides?mobile=1234567890
	mobile := r.URL.Query().Get("mobile")
	if mobile == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Mobile number is required"}`))
		return
	}

	// 2. Query with JOIN to filter by phone number and status
	// We use 'u' for usertbl and 'r' for RideBookingstbl
	query := `
		SELECT 
			r.BookingID, r.service_type, r.customer_name, r.pickup_location, r.drop_location, 
			r.address, r.start_date, r.end_date, r.start_time, r.no_of_days, r.no_of_hours, 
			r.has_car, r.car_name, r.car_number, r.total_fare, r.distance, r.status 
		FROM RideBookingstbl r
		JOIN usertbl u ON r.UserID = u.UserID
		WHERE u.phoneNumber = ? AND r.status = 'Completed'
		ORDER BY r.updated_at DESC`

	rows, err := db.DB.Query(query, mobile)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Database query failed"}`))
		return
	}
	defer rows.Close()

	var rides []RideResponse // Using your existing RideResponse struct
	for rows.Next() {
		var ride RideResponse
		err := rows.Scan(
			&ride.BookingID, &ride.ServiceType, &ride.CustomerName, &ride.PickupLocation,
			&ride.DropLocation, &ride.Address, &ride.StartDate, &ride.EndDate,
			&ride.StartTime, &ride.NoOfDays, &ride.NoOfHours, &ride.HasCar,
			&ride.CarName, &ride.CarNumber, &ride.TotalFare, &ride.Distance, &ride.Status,
		)
		if err != nil {
			continue
		}
		rides = append(rides, ride)
	}

	// Ensure we return [] instead of null if no history exists
	if len(rides) == 0 {
		rides = []RideResponse{}
	}

	json.NewEncoder(w).Encode(rides)
}
