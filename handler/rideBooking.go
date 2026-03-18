package handlers

import (
	"database/sql"
	db "dod-backend/Database"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// 1. Define the Struct to match the Angular JSON payload
type BookingRequest struct {
	ServiceType    string   `json:"serviceType"`
	PickupLocation *string  `json:"pickupLocation"`
	DropLocation   *string  `json:"dropLocation"`
	Address        *string  `json:"address"`
	StartDate      *string  `json:"startDate"`
	EndDate        *string  `json:"endDate"`
	StartTime      *string  `json:"startTime"`
	NoOfDays       *int     `json:"noOfDays"`
	NoOfHours      *int     `json:"noOfHours"`
	HasCar         *bool    `json:"hasCar"`
	CarName        *string  `json:"carName"`
	CarNumber      *string  `json:"carNumber"`
	TotalFare      *float64 `json:"totalFare"`
	Distance       *float64 `json:"distance"`
	CustomerName   *string  `json:"customerName"` // Getting this from frontend now
}

// 2. The Handler Function
func BookRideHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// 1. Extract the phoneNumber from the URL
	pathParts := strings.Split(r.URL.Path, "/")
	phoneNumber := pathParts[len(pathParts)-1]

	if phoneNumber == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Phone number is required"}`))
		return
	}

	// 2. Find the UserID associated with this phone number
	var userID int
	// We only need the UserID now since Name comes from the frontend
	err := db.DB.QueryRow("SELECT UserID FROM usertbl WHERE phoneNumber = ?", phoneNumber).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"error": "User not found"}`))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "Database error checking user"}`))
		}
		return
	}

	// 3. Decode the incoming JSON body
	var req BookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Invalid JSON format"}`))
		return
	}

	// 4. Insert the booking into the RideBookingstbl
	// Added 'DriverID' to the column list and the VALUES placeholders
	query := `INSERT INTO RideBookingstbl
    (UserID, DriverID, customer_name, service_type, pickup_location, drop_location, address, 
     start_date, end_date, start_time, no_of_days, no_of_hours, 
     has_car, car_name, car_number, total_fare, distance) 
    VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	// 5. Execute the query
	// Included the hardcoded DriverID (1) as the second argument
	_, err = db.DB.Exec(query,
		userID,
		1, // Hardcoded DriverID for 'Abhi'
		req.CustomerName,
		req.ServiceType,
		req.PickupLocation,
		req.DropLocation,
		req.Address,
		req.StartDate,
		req.EndDate,
		req.StartTime,
		req.NoOfDays,
		req.NoOfHours,
		req.HasCar,
		req.CarName,
		req.CarNumber,
		req.TotalFare,
		req.Distance,
	)

	if err != nil {
		fmt.Println("Database Insert Error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Failed to save booking"}`))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "Ride booked successfully!"}`))
}
