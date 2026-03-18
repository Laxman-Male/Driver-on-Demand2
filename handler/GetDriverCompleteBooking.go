package handlers

import (
	db "dod-backend/Database"
	"encoding/json"
	"fmt"
	"net/http"
)

// GetDriverCompletedRides fetches history for a specific driver based on mobile

func GetDriverCompletedRides(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	mobile := r.URL.Query().Get("mobile")
	if mobile == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Driver mobile is required"}`))
		return
	}

	// Explicitly selecting all columns to match your struct
	query := `
        SELECT 
            r.BookingID, r.service_type, r.customer_name, r.pickup_location, r.drop_location, 
            r.address, r.start_date, r.end_date, r.start_time, r.no_of_days, r.no_of_hours, 
            r.has_car, r.car_name, r.car_number, r.total_fare, r.distance, r.status 
        FROM RideBookingstbl r
        JOIN DriverDetailstbl d ON r.DriverID = d.DriverID
        WHERE d.PhoneNumber = ? AND r.status = 'Completed'
        ORDER BY r.updated_at DESC`

	rows, err := db.DB.Query(query, mobile)
	if err != nil {
		fmt.Println("DB Query Error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	rides := []RideResponse{} // Initialize as empty slice to avoid 'null' in JSON

	for rows.Next() {
		var ride RideResponse
		err := rows.Scan(
			&ride.BookingID, &ride.ServiceType, &ride.CustomerName, &ride.PickupLocation,
			&ride.DropLocation, &ride.Address, &ride.StartDate, &ride.EndDate,
			&ride.StartTime, &ride.NoOfDays, &ride.NoOfHours, &ride.HasCar,
			&ride.CarName, &ride.CarNumber, &ride.TotalFare, &ride.Distance, &ride.Status,
		)

		if err != nil {
			// THIS IS KEY: It will tell you why your array is blank
			fmt.Println("Scan Error on BookingID:", ride.BookingID, " | Error:", err)
			continue
		}
		rides = append(rides, ride)
	}

	json.NewEncoder(w).Encode(rides)
}
