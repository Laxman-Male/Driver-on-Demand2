package handlers

import (
	db "dod-backend/Database"
	"encoding/json"
	"fmt"
	"net/http"
	// db "dod-backend/Database" // Ensure your DB package is imported
)

// 1. Struct for sending data to the frontend
// We use pointers (*string, *int) so Go doesn't crash when it hits a NULL in the database!
type RideResponse struct {
	BookingID      int      `json:"bookingId"`
	ServiceType    string   `json:"serviceType"`
	PickupLocation *string  `json:"pickupLocation"`
	DropLocation   *string  `json:"dropLocation"`
	Address        *string  `json:"address"`
	StartDate      *string  `json:"startDate"` // Dates come out as strings in Go JSON
	EndDate        *string  `json:"endDate"`
	StartTime      *string  `json:"startTime"`
	NoOfDays       *int     `json:"noOfDays"`
	NoOfHours      *int     `json:"noOfHours"`
	HasCar         *bool    `json:"hasCar"`
	CarName        *string  `json:"carName"`
	CarNumber      *string  `json:"carNumber"`
	TotalFare      *float64 `json:"totalFare"`
	Distance       *float64 `json:"distance"`
	Status         string   `json:"status"`
	CustomerName   string   `json:"customerName"`
}

// 2. The Handler Function
func GetPendingRidesForDriver(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// 1. Write the SQL Query to fetch only 'Pending' rides, newest first
	query := `SELECT BookingID,customer_name, service_type, pickup_location, drop_location, address, 
			  start_date, end_date, start_time, no_of_days, no_of_hours, 
			  has_car, car_name, car_number, total_fare, distance, status 
			  FROM RideBookingstbl 
			  WHERE status = 'Pending' 
			  ORDER BY created_at DESC`

	// 2. Execute the query using db.Query (not QueryRow, because we expect multiple rides!)
	rows, err := db.DB.Query(query)
	if err != nil {
		fmt.Println("--", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Failed to fetch rides from database"}`))
		return
	}
	defer rows.Close() // ALWAYS close your rows to prevent memory leaks!

	// 3. Create a slice (array) to hold all the rides
	var rides []RideResponse

	// 4. Loop through the results row by row
	for rows.Next() {
		var ride RideResponse

		// Scan the row columns into our struct
		// The order here MUST perfectly match the order in your SELECT statement above!
		err := rows.Scan(
			&ride.BookingID, &ride.CustomerName, &ride.ServiceType, &ride.PickupLocation, &ride.DropLocation, &ride.Address,
			&ride.StartDate, &ride.EndDate, &ride.StartTime, &ride.NoOfDays, &ride.NoOfHours,
			&ride.HasCar, &ride.CarName, &ride.CarNumber, &ride.TotalFare, &ride.Distance, &ride.Status,
		)

		if err != nil {
			continue // If one row fails to scan, skip it and keep going
		}

		// Add the successfully scanned ride to our slice
		rides = append(rides, ride)
	}

	// 5. Ensure we return an empty JSON array [] instead of 'null' if there are no rides
	if len(rides) == 0 {
		rides = []RideResponse{}
	}

	fmt.Println("---", rides)
	// 6. Send the array of rides back to Angular!
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(rides)
}
