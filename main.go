package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "github.com/google/uuid"
    "strings"
)

var receipts = make(map[string]int)

type Receipt struct {
    Retailer     string  json:"retailer"
    Total        float64 json:"total"
    PurchaseDate string  json:"purchaseDate"
    PurchaseTime string  json:"purchaseTime"
    Items        []Item  json:"items"
}

type Item struct {
    Description string  json:"description"
    Price       float64 json:"price"
}

type ProcessResponse struct {
    ID string json:"id"
}

type PointsResponse struct {
    Points int json:"points"
}

func main() {
    // Set up the HTTP routes
    http.HandleFunc("/receipts/process", processReceiptHandler)
    http.HandleFunc("/receipts/", getPointsHandler)

    // Log that the server is starting
    fmt.Println("Server is running on port 8080...")

    // Start the HTTP server and block execution
    log.Fatal(http.ListenAndServe(":8080", nil))
}

// Handler for processing receipts
func processReceiptHandler(w http.ResponseWriter, r *http.Request) {
    // Only accept POST requests
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
        return
    }

    var receipt Receipt
    // Decode JSON from the request body
    if err := json.NewDecoder(r.Body).Decode(&receipt); err != nil {
        fmt.Println("Error decoding JSON:", err)  // Debugging line
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }

    // Generate a unique ID for the receipt
    id := uuid.New().String()

    // Calculate points for the receipt
    points := calculatePoints(receipt)

    // Store the points with the generated ID
    receipts[id] = points

    // Send the ID as the response
    response := ProcessResponse{ID: id}
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

// Handler for retrieving points for a receipt by ID
func getPointsHandler(w http.ResponseWriter, r *http.Request) {
    // Ensure it's a GET request
    if r.Method != http.MethodGet {
        http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
        return
    }

    // Extract the ID and check if URL contains "/points" at the end
    pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/receipts/"), "/")
    if len(pathParts) != 2 || pathParts[1] != "points" {
        http.Error(w, "Invalid endpoint", http.StatusNotFound)
        return
    }
    id := pathParts[0]

    // Look up points by ID
    points, found := receipts[id]
    if !found {
        http.Error(w, "Receipt not found", http.StatusNotFound)
        return
    }

    // Send the points as the response
    response := PointsResponse{Points: points}
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

// Function to calculate the points for a receipt
func calculatePoints(receipt Receipt) int {
    points := 0

    // 1. One point for every alphanumeric character in the retailer name
    for _, c := range receipt.Retailer {
        if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') {
            points++
        }
    }

    // 2. 50 points if the total is a round dollar amount with no cents
    if receipt.Total == float64(int(receipt.Total)) {
        points += 50
    }

    // 3. 25 points if the total is a multiple of 0.25
    if int(receipt.Total*4)%4 == 0 {
        points += 25
    }

    // 4. 5 points for every two items on the receipt
    points += (len(receipt.Items) / 2) * 5

    // 5. For each item, if the description length is a multiple of 3, multiply the price by 0.2 and round up
    for _, item := range receipt.Items {
        if len(item.Description)%3 == 0 {
            points += int(item.Price*0.2 + 0.5) // rounding up
        }
    }

    // 6. 6 points if the day in the purchase date is odd
    var year, month, day int
    fmt.Sscanf(receipt.PurchaseDate, "%d-%d-%d", &year, &month, &day)
    if day%2 != 0 {
        points += 6
    }

    // 7. 10 points if the time is after 2:00 PM and before 4:00 PM
    var hour, minute int
    fmt.Sscanf(receipt.PurchaseTime, "%d:%d", &hour, &minute)
    if hour >= 14 && hour < 16 {
        points += 10
    }

    return points
}






// package main

// import (
//     "encoding/json"
//     "fmt"
//     "log"
//     "net/http"
//     "github.com/google/uuid"
// )

// var receipts = make(map[string]int)

// type Receipt struct {
//     Retailer     string  json:"retailer"
//     Total        float64 json:"total"
//     PurchaseDate string  json:"purchaseDate"
//     PurchaseTime string  json:"purchaseTime"
//     Items        []Item  json:"items"
// }

// type Item struct {
//     Description string  json:"description"
//     Price       float64 json:"price"
// }

// type ProcessResponse struct {
//     ID string json:"id"
// }

// type PointsResponse struct {
//     Points int json:"points"
// }

// func main() {
//     // Set up the HTTP routes
//     http.HandleFunc("/receipts/process", processReceiptHandler)
//     http.HandleFunc("/receipts/", getPointsHandler)

//     // Log that the server is starting
//     fmt.Println("Server is running on port 8080...")

//     // Start the HTTP server and block execution
//     log.Fatal(http.ListenAndServe(":8080", nil))
// }

// // Handler for processing receipts
// func processReceiptHandler(w http.ResponseWriter, r *http.Request) {
//     // Only accept POST requests
//     if r.Method != http.MethodPost {
//         http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
//         return
//     }

//     var receipt Receipt
//     // Decode JSON from the request body
//     if err := json.NewDecoder(r.Body).Decode(&receipt); err != nil {
//         fmt.Println("Error decoding JSON:", err)  // Debugging line
//         http.Error(w, "Invalid input", http.StatusBadRequest)
//         return
//     }

//     // Generate a unique ID for the receipt
//     id := uuid.New().String()

//     // Calculate points for the receipt
//     points := calculatePoints(receipt)

//     // Store the points with the generated ID
//     receipts[id] = points

//     // Send the ID as the response
//     response := ProcessResponse{ID: id}
//     w.Header().Set("Content-Type", "application/json")
//     json.NewEncoder(w).Encode(response)
// }

// // Handler for retrieving points for a receipt by ID
// func getPointsHandler(w http.ResponseWriter, r *http.Request) {
//     // Ensure it's a GET request
//     if r.Method != http.MethodGet {
//         http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
//         return
//     }

//     id := r.URL.Path[len("/receipts/"):] // Get the ID from the URL path
//     points, found := receipts[id]
//     if !found {
//         http.Error(w, "Receipt not found", http.StatusNotFound)
//         return
//     }

//     // Send the points as the response
//     response := PointsResponse{Points: points}
//     w.Header().Set("Content-Type", "application/json")
//     json.NewEncoder(w).Encode(response)
// }

// // Function to calculate the points for a receipt
// func calculatePoints(receipt Receipt) int {
//     points := 0

//     // 1. One point for every alphanumeric character in the retailer name
//     for _, c := range receipt.Retailer {
//         if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') {
//             points++
//         }
//     }

//     // 2. 50 points if the total is a round dollar amount with no cents
//     if receipt.Total == float64(int(receipt.Total)) {
//         points += 50
//     }

//     // 3. 25 points if the total is a multiple of 0.25
//     if int(receipt.Total*4)%4 == 0 {
//         points += 25
//     }

//     // 4. 5 points for every two items on the receipt
//     points += (len(receipt.Items) / 2) * 5

//     // 5. For each item, if the description length is a multiple of 3, multiply the price by 0.2 and round up
//     for _, item := range receipt.Items {
//         if len(item.Description)%3 == 0 {
//             points += int(item.Price*0.2 + 0.5) // rounding up
//         }
//     }

//     // 6. 6 points if the day in the purchase date is odd
//     var year, month, day int
//     fmt.Sscanf(receipt.PurchaseDate, "%d-%d-%d", &year, &month, &day)
//     if day%2 != 0 {
//         points += 6
//     }

//     // 7. 10 points if the time is after 2:00 PM and before 4:00 PM
//     var hour, minute int
//     fmt.Sscanf(receipt.PurchaseTime, "%d:%d", &hour, &minute)
//     if hour >= 14 && hour < 16 {
//         points += 10
//     }

//     return points
// }