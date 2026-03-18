package handlers

import (
	db "dod-backend/Database"
	"encoding/json"
	"net/http"
)

func GetAllUserBookings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// 1. Get mobile from query param: /get-all-bookings?mobile=9876543210
	mobile := r.URL.Query().Get("mobile")
	if mobile == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Mobile number is required"}`))
		return
	}

	// 2. SQL Query with JOIN
	// We fetch every ride where the user's phone number matches
	query := `
		SELECT 
			r.BookingID, r.service_type, r.customer_name, r.pickup_location, r.drop_location, 
			r.address, r.start_date, r.end_date, r.start_time, r.no_of_days, r.no_of_hours, 
			r.has_car, r.car_name, r.car_number, r.total_fare, r.distance, r.status 
		FROM RideBookingstbl r
		JOIN usertbl u ON r.UserID = u.UserID
		WHERE u.phoneNumber = ?
		ORDER BY r.created_at DESC`

	rows, err := db.DB.Query(query, mobile)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Database error fetching bookings"}`))
		return
	}
	defer rows.Close()

	var rides []RideResponse
	for rows.Next() {
		var r RideResponse
		err := rows.Scan(
			&r.BookingID, &r.ServiceType, &r.CustomerName, &r.PickupLocation,
			&r.DropLocation, &r.Address, &r.StartDate, &r.EndDate,
			&r.StartTime, &r.NoOfDays, &r.NoOfHours, &r.HasCar,
			&r.CarName, &r.CarNumber, &r.TotalFare, &r.Distance, &r.Status,
		)
		if err != nil {
			continue // Skip rows with errors
		}
		rides = append(rides, r)
	}

	// 3. Handle empty case
	if len(rides) == 0 {
		rides = []RideResponse{}
	}

	json.NewEncoder(w).Encode(rides)
}
