package main

import (
	"database/sql"
	"fmt"
	"log"
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

func main() {
	// Open a connection to the database
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/ddb")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Insert a new medication
	medication := Medication{Name: "Medicine A", Dosage: "10mg", Manufacturer: "Manufacturer A", Price: 10.50}
	medicationID, err := insertMedication(db, medication)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted medication with ID:", medicationID)

	// Update the medication
	medication.Price = 11.50
	err = updateMedication(db, medication)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Updated medication")

	// Delete the medication
	err = deleteMedication(db, medicationID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Deleted medication")

	// Query all medications
	medications, err := queryAllMedications(db)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("All medications:", medications)
	// Get a single medication
	medID := 1
	medication, err = getMedication(db, medID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Medication with ID", medID, ":", medication)

	// Get all inventory items
	inventoryItems, err := queryAllInventory(db)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("All inventory items:", inventoryItems)

	// Get all prescriptions
	prescriptions, err := queryAllPrescriptions(db)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("All prescriptions:", prescriptions)

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

func queryAllInventory(db *sql.DB) ([]Inventory, error) {
	rows, err := db.Query("SELECT * FROM inventory")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	inventoryItems := []Inventory{}
	for rows.Next() {
		var inventoryItem Inventory
		err := rows.Scan(&inventoryItem.MedicationID, &inventoryItem.Quantity)
		if err != nil {
			return nil, err
		}
		inventoryItems = append(inventoryItems, inventoryItem)
	}
	return inventoryItems, nil
}

func queryAllPrescriptions(db *sql.DB) ([]Prescription, error) {
	rows, err := db.Query("SELECT * FROM prescriptions")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	prescriptions := []Prescription{}
	for rows.Next() {
		var prescription Prescription
		err := rows.Scan(&prescription.ID, &prescription.MedicationID, &prescription.Quantity, &prescription.PatientName, &prescription.DoctorName, &prescription.DatePrescribed)
		if err != nil {
			return nil, err
		}
		prescriptions = append(prescriptions, prescription)
	}
	return prescriptions, nil
}
