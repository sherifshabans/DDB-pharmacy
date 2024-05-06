package main

import (
	"fmt"
	"net"
	"strconv"
)

type Medication struct {
	ID           int
	Name         string
	Dosage       string
	Manufacturer string
	Price        float64
}

func main() {
	// Connect to server
	conn, err := net.Dial("tcp", "192.168.1.7:8080")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	//createSendToServer(conn)
	getAllSendToServer(conn)
}

func createSendToServer(conn net.Conn) {
	// Send request to server
	request := []byte("")
	_, _ = conn.Write(request)
	med := Medication{ID: 1, Name: "Ahmed", Dosage: "Abdo", Manufacturer: "Hussein", Price: 15.5}
	create(conn, med)
}

func getAllSendToServer(conn net.Conn) {
	getAll(conn)
	get := make([]byte, 1024)
	conn.Read(get)
	fmt.Println(get)
}

func create(conn net.Conn, medication Medication) {
	request := []byte("insert ")

	x := []byte(strconv.Itoa(medication.ID))
	x = append(x, " "...)
	request = append(request, x...)

	x = []byte(medication.Name)
	x = append(x, " "...)
	request = append(request, x...)

	x = []byte(medication.Dosage)
	x = append(x, " "...)
	request = append(request, x...)

	x = []byte(medication.Manufacturer)
	x = append(x, " "...)
	request = append(request, x...)

	x = []byte(strconv.FormatFloat(medication.Price, 'f', -1, 64))
	request = append(request, x...)

	conn.Write(request)
}

// Function to read all medications from the database
func getAll(conn net.Conn) {
	request := []byte("getAll")
	conn.Write(request)
}
