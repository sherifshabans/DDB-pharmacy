package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type Department struct {
	id             int
	DepartmentName string
}

type Student struct {
	id           int
	StudentName  string
	DepartmentID int
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
	case "insertDepartment":
		if parts[0] != "insertDepartment" {
			fmt.Println("Invalid request format for insert Department")
			return
		}

		DepartmentName := parts[2]

		department := Department{DepartmentName: DepartmentName}

		departmentID, err := insertDepartment(db, department)
		if err != nil {
			fmt.Println("Error inserting department:", err)
			return
		}
		fmt.Println("Inserted department with ID:", departmentID)
	case "queryAllDepartments":
		if parts[0] != "queryAllDepartments" {
			fmt.Println("Invalid request format on query all departments")
			return
		}
		departments, err := queryAllDepartments(db)
		if err != nil {
			fmt.Println("Error querying all departments:", err)
			return
		}

		var data []byte
		for _, dep := range departments {
			data = append(data, []byte(fmt.Sprintf("%d|%s\n", dep.id, dep.DepartmentName))...)
		}

		_, err = conn.Write(data)
		if err != nil {
			fmt.Println("Error sending data to client on query All Departments:", err)
			return
		}
		fmt.Println("Sent to client successfully")
	case "updateDepartment":
		if parts[0] != "updateDepartment" {
			fmt.Println("Invalid request format on update department")
			return
		}

		idStr := parts[1]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			fmt.Println("Error parsing ID on update department:", err)
			return
		}
		DepartmentName := parts[2]

		department := Department{id: id, DepartmentName: DepartmentName}

		err = updateDepartment(db, department)
		if err != nil {
			fmt.Println("Error updating department:", err)
			return
		}
		message := "Updated successfully"
		bytes := []byte(message)
		// Send byte array back to client
		_, err = conn.Write(bytes)
		if err != nil {
			fmt.Println("Error sending department data to client:", err)
			return
		}
		fmt.Println("Sent to client and Updated successfully")
	case "deleteDepartment":
		if parts[0] != "deleteDepartment" {
			fmt.Println("Invalid request format for delete department")
			return
		}

		idStr := parts[1]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			fmt.Println("Error parsing department ID:", err)
			return
		}

		err = deleteDepartment(db, id)
		if err != nil {
			fmt.Println("Error deleting department:", err)
			return
		}
		message := "Deleted successfully"
		bytes := []byte(message)
		// Send byte array back to client
		_, err = conn.Write(bytes)
		if err != nil {
			fmt.Println("Error sending message to client:", err)
			return
		}
		fmt.Println("Sent to client and Deleted successfully")
	case "getDepartment":
		if parts[0] != "getDepartment" {
			fmt.Println("Invalid request format for get department")
			return
		}

		idStr := parts[1]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			fmt.Println("Error parsing department ID:", err)
			return
		}

		department, err := getDepartment(db, id)
		if err != nil {
			fmt.Println("Error getting department:", err)
			return
		}

		data := []byte(fmt.Sprintf("%d|%s\n", department.id, department.DepartmentName))

		// Send byte array back to client
		_, err = conn.Write(data)
		if err != nil {
			fmt.Println("Error sending department data to client:", err)
			return
		}
		fmt.Println("Sent to client successfully")
	case "insertStudent":
		if parts[0] != "insertStudent" {
			fmt.Println("Invalid request format for insert Student")
			return
		}
		studentName := parts[1]
		departmentIDStr := parts[2]
		departmentID, err := strconv.Atoi(departmentIDStr)
		if err != nil {
			fmt.Println("Error converting departmentID to int:", err)
			return
		}

		student := Student{StudentName: studentName, DepartmentID: departmentID}

		studentID, err := insertStudent(db, student)
		if err != nil {
			fmt.Println("Error inserting student:", err)
			return
		}
		fmt.Println("Inserted student with ID:", studentID)
	case "queryAllStudents":
		if parts[0] != "queryAllStudents" {
			fmt.Println("Invalid request format on query all students")
			return
		}
		students, err := queryAllStudents(db)
		if err != nil {
			fmt.Println("Error querying all students:", err)
			return
		}

		var data []byte
		for _, stu := range students {
			data = append(data, []byte(fmt.Sprintf("%d|%s|%d\n", stu.id, stu.StudentName, stu.DepartmentID))...)
		}

		_, err = conn.Write(data)
		if err != nil {
			fmt.Println("Error sending data to client on query All students:", err)
			return
		}
		fmt.Println("Sent to client successfully")
	case "updateStudent":
		if parts[0] != "updateStudent" {
			fmt.Println("Invalid request format on update student")
			return
		}

		idStr := parts[1]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			fmt.Println("Error parsing ID on update department:", err)
			return
		}
		StudentName := parts[2]
		departmentIDStr := parts[3]
		departmentID, err := strconv.Atoi(departmentIDStr)
		if err != nil {
			fmt.Println("Error parsing ID on update department:", err)
			return
		}

		student := Student{id: id, StudentName: StudentName, DepartmentID: departmentID}

		err = updateStudent(db, student)
		if err != nil {
			fmt.Println("Error updating department:", err)
			return
		}
		message := "Updated successfully"
		bytes := []byte(message)
		// Send byte array back to client
		_, err = conn.Write(bytes)
		if err != nil {
			fmt.Println("Error sending message to client:", err)
			return
		}
		fmt.Println("Sent to client and Updated successfully")
	case "deleteStudent":
		if parts[0] != "deleteStudent" {
			fmt.Println("Invalid request format for delete student")
			return
		}

		idStr := parts[1]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			fmt.Println("Error parsing student ID:", err)
			return
		}

		err = deleteStudent(db, id)
		if err != nil {
			fmt.Println("Error deleting student:", err)
			return
		}
		message := "Deleted successfully"
		bytes := []byte(message)
		// Send byte array back to client
		_, err = conn.Write(bytes)
		if err != nil {
			fmt.Println("Error sending message to client:", err)
			return
		}
		fmt.Println("Sent to client and Deleted successfully")
	case "getStudent":
		if parts[0] != "getStudent" {
			fmt.Println("Invalid request format for get student")
			return
		}

		idStr := parts[1]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			fmt.Println("Error parsing department ID:", err)
			return
		}

		student, err := getStudent(db, id)
		if err != nil {
			fmt.Println("Error getting student:", err)
			return
		}

		data := []byte(fmt.Sprintf("%d|%s\n", student.id, student.StudentName, student.DepartmentID))

		// Send byte array back to client
		_, err = conn.Write(data)
		if err != nil {
			fmt.Println("Error sending student data to client:", err)
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

func dbConfig() (*sql.DB, error) {
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
	db, err := dbConfig()
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

func insertDepartment(db *sql.DB, m Department) (int, error) {
	result, err := db.Exec("INSERT INTO Department (DepartmentName) VALUES (?)", m.DepartmentName)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func queryAllDepartments(db *sql.DB) ([]Department, error) {
	rows, err := db.Query("SELECT * FROM Department")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	department := []Department{}
	for rows.Next() {
		var m Department
		err := rows.Scan(&m.id, &m.DepartmentName)
		if err != nil {
			return nil, err
		}
		department = append(department, m)
	}
	return department, nil
}

func updateDepartment(db *sql.DB, m Department) error {
	_, err := db.Exec("UPDATE Department SET DepartmentName=? WHERE id=?", m.DepartmentName, m.id)
	return err
}

func deleteDepartment(db *sql.DB, id int) error {
	_, err := db.Exec("DELETE FROM Department WHERE id=?", id)
	return err
}

func getDepartment(db *sql.DB, id int) (Department, error) {
	var department Department
	err := db.QueryRow("SELECT * FROM Department WHERE id = ?", id).Scan(&department.DepartmentName)
	if err != nil {
		return Department{}, err
	}
	return department, nil
}

// Student

func insertStudent(db *sql.DB, m Student) (int, error) {
	result, err := db.Exec("INSERT INTO Student (StudentName,DepartmentID) VALUES (?,?)", m.StudentName, m.DepartmentID)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func updateStudent(db *sql.DB, m Student) error {
	_, err := db.Exec("UPDATE Student SET StudentName=?, DepartmentID=? WHERE id = ?", m.StudentName, m.DepartmentID, m.id)
	return err
}

func deleteStudent(db *sql.DB, id int) error {
	_, err := db.Exec("DELETE FROM Student WHERE id=?", id)
	return err
}

func queryAllStudents(db *sql.DB) ([]Student, error) {
	rows, err := db.Query("SELECT * FROM Student")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	student := []Student{}
	for rows.Next() {
		var m Student
		err := rows.Scan(&m.id, &m.StudentName, &m.DepartmentID)
		if err != nil {
			return nil, err
		}
		student = append(student, m)
	}
	return student, nil
}

func getStudent(db *sql.DB, id int) (Student, error) {
	var student Student
	err := db.QueryRow("SELECT * FROM Student WHERE id = ?", id).Scan(&student.StudentName, &student.DepartmentID)
	if err != nil {
		return Student{}, err
	}
	return student, nil
}
