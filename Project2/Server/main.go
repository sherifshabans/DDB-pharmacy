package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Medication struct {
	ID           int
	Name         string
	Dosage       string
	Manufacturer string
	Price        float64
}

type Inventory struct {
	MedicationID int
	Quantity     int
}

type Prescription struct {
	ID             int
	MedicationID   int
	Quantity       int
	PatientName    string
	DoctorName     string
	DatePrescribed time.Time
}

func handleRequest(conn net.Conn, db *sql.DB) {
	defer conn.Close()

	// Read request from client
	request := make([]byte, 1024)
	n, err := conn.Read(request)
	if err != nil {
		fmt.Println("Error reading request:", err)
		return
	}

	requestStr := string(request[:n])
	parts := strings.Split(requestStr, " ")
	switchStr := parts[0]

	// Process request and execute server-side function
	switch string(switchStr) {
	case "insert":
		if len(parts) < 5 || parts[0] != "insert" {
			fmt.Println("Invalid request format")
			return
		}
		// Extract medication data from request
		id, err := strconv.Atoi(parts[1])
		if err != nil {
			fmt.Println("Error parsing id:", err)
			return
		}
		name := parts[2]
		dosage := parts[3]
		manufacturer := parts[4]
		priceStr := parts[5]
		price, err := strconv.ParseFloat(priceStr, 64)
		if err != nil {
			fmt.Println("Error parsing price:", err)
			return
		}

		// Create medication object
		medication := Medication{ID: id, Name: name, Dosage: dosage, Manufacturer: manufacturer, Price: price}

		// Insert medication into database
		medicationID, err := insertMedication(db, medication)
		if err != nil {
			fmt.Println("Error inserting medication:", err)
			return
		}
		fmt.Println("Inserted medication with ID:", medicationID)
	case "getAll":
		medications, err := queryAllMedications(db)
		if err != nil {
			fmt.Println("Error querying all medications:", err)
			return
		}
		// Convert medications to a byte slice
		var data []byte
		for _, med := range medications {
			data = append(data, []byte(fmt.Sprintf("%d|%s|%s|%s|%.2f\n", med.ID, med.Name, med.Dosage, med.Manufacturer, med.Price))...)
		}

		// Send byte array back to client
		_, err = conn.Write(data)
		if err != nil {
			fmt.Println("Error sending data to client:", err)
			return
		}
		fmt.Println("Sent to client successfully")
	case "update":
		if len(parts) < 6 || parts[0] != "update" {
			fmt.Println("Invalid request format")
			return
		}
		// Extract medication data from request
		idStr := parts[1]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			fmt.Println("Error parsing ID:", err)
			return
		}
		name := parts[2]
		dosage := parts[3]
		manufacturer := parts[4]
		priceStr := parts[5]
		price, err := strconv.ParseFloat(priceStr, 64)
		if err != nil {
			fmt.Println("Error parsing price:", err)
			return
		}

		// Create medication object
		medication := Medication{ID: id, Name: name, Dosage: dosage, Manufacturer: manufacturer, Price: price}

		// Update medication in database
		err = updateMedication(db, medication)
		if err != nil {
			fmt.Println("Error updating medication:", err)
			return
		}
		message := "Updated successfully"
		bytes := []byte(message)
		// Send byte array back to client
		_, err = conn.Write(bytes)
		if err != nil {
			fmt.Println("Error sending data to client:", err)
			return
		}
		fmt.Println("Sent to client and Updated successfully")
	case "delete":
		if len(parts) < 2 || parts[0] != "delete" {
			fmt.Println("Invalid request format")
			return
		}
		// Extract medication ID from request
		idStr := parts[1]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			fmt.Println("Error parsing ID:", err)
			return
		}

		// Delete medication from database
		err = deleteMedication(db, id)
		if err != nil {
			fmt.Println("Error deleting medication:", err)
			return
		}
		message := "Deleted successfully"
		bytes := []byte(message)
		// Send byte array back to client
		_, err = conn.Write(bytes)
		if err != nil {
			fmt.Println("Error sending data to client:", err)
			return
		}
		fmt.Println("Sent to client and Deleted successfully")
	case "getOne":
		if len(parts) < 2 || parts[0] != "getOne" {
			fmt.Println("Invalid request format")
			return
		}
		// Extract medication ID from request
		idStr := parts[1]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			fmt.Println("Error parsing ID:", err)
			return
		}

		// Get medication from database
		medication, err := getMedication(db, id)
		if err != nil {
			fmt.Println("Error getting medication:", err)
			return
		}

		// Convert medication to a byte slice
		data := []byte(fmt.Sprintf("%d|%s|%s|%s|%.2f\n", medication.ID, medication.Name, medication.Dosage, medication.Manufacturer, medication.Price))

		// Send byte array back to client
		_, err = conn.Write(data)
		if err != nil {
			fmt.Println("Error sending data to client:", err)
			return
		}
		fmt.Println("Sent to client successfully")

	default:
		fmt.Println("Invalid request")
	}
}

func createDatabaseSchema(db *sql.DB) error {
	// Read the SQL file
	sqlFile, err := ioutil.ReadFile("db.sql")
	if err != nil {
		return err
	}

	// Split the file contents into individual queries
	queries := strings.Split(string(sqlFile), ";")

	// Execute each query to create the database schema
	for _, query := range queries {
		trimmedQuery := strings.TrimSpace(query)
		if trimmedQuery != "" {
			_, err := db.Exec(trimmedQuery)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func dbconfig() (*sql.DB, error) {
	// Open a connection to the database
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	// Create the database if it doesn't exist
	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS ddb")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	// Switch to the database
	_, err = db.Exec("USE ddb")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	// Create the database schema from the SQL file
	err = createDatabaseSchema(db)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return db, nil
}

func main() {

	// Database config
	db, err := dbconfig()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ip := "0.0.0.0:8080"
	// Start server
	listener, err := net.Listen("tcp", ip) // Listen for incoming connections on specified address and port
	if err != nil {
		fmt.Println("Error listening:", err.Error()) // Print error if unable to start listening
		return
	}
	defer listener.Close()                               // Close the listener when main function returns
	fmt.Println("Server started, listening on port", ip) // Print a message indicating server startup

	// Accept incoming connections
	for {
		connection, err := listener.Accept() // Accept incoming connection
		if err != nil {
			fmt.Println("Error accepting connection:", err.Error()) // Print error if unable to accept connection
			continue
		}

		// handle request coming from the client
		go handleRequest(connection, db) // Handle the connection in a separate goroutine
	}

}

func insertMedication(db *sql.DB, m Medication) (int, error) {
	result, err := db.Exec("INSERT INTO medications (name, dosage, manufacturer, price) VALUES (?, ?, ?, ?)", m.Name, m.Dosage, m.Manufacturer, m.Price)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func updateMedication(db *sql.DB, m Medication) error {
	_, err := db.Exec("UPDATE medications SET name=?, dosage=?, manufacturer=?, price=? WHERE id=?", m.Name, m.Dosage, m.Manufacturer, m.Price, m.ID)
	return err
}

func deleteMedication(db *sql.DB, id int) error {
	_, err := db.Exec("DELETE FROM medications WHERE id=?", id)
	return err
}

func queryAllMedications(db *sql.DB) ([]Medication, error) {
	rows, err := db.Query("SELECT * FROM medications")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	medications := []Medication{}
	for rows.Next() {
		var m Medication
		err := rows.Scan(&m.ID, &m.Name, &m.Dosage, &m.Manufacturer, &m.Price)
		if err != nil {
			return nil, err
		}
		medications = append(medications, m)
	}
	return medications, nil
}

func getMedication(db *sql.DB, id int) (Medication, error) {
	var medication Medication
	err := db.QueryRow("SELECT * FROM medications WHERE id = ?", id).Scan(&medication.ID, &medication.Name, &medication.Dosage, &medication.Manufacturer, &medication.Price)
	if err != nil {
		return Medication{}, err
	}
	return medication, nil
}
