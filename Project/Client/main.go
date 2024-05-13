package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
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

func main() {
	for {
		// Connect to server
		conn, err := net.Dial("tcp", "192.168.62.14:8080")
		if err != nil {
			fmt.Println("Error connecting to server:", err)
			return
		}
		defer conn.Close()

		fmt.Print("Enter command (getAllDepartments, getAllStudents, createDepartment, createStudent, deleteDepartment,  deleteStudent, updateDepartment, updateStudent): ")
		var command string
		fmt.Scanln(&command)

		switch command {
		case "getAllDepartments":
			getAllDepartments(conn)
		case "createDepartment":
			createSendToServerDep(conn)
		case "deleteDepartment":
			deleteDepartment(conn)
		case "updateDepartment":
			updateDepartment(conn)
		case "getAllStudents":
			getAllStudents(conn)
		case "createStudent":
			createSendToServerStud(conn)
		case "deleteStudent":
			deleteStudent(conn)
		case "updateStudent":
			updateStudent(conn)
		case "exit":
			os.Exit(0)
		default:
			fmt.Println("Invalid command. Please try again.")
		}
	}
}

func createSendToServerDep(conn net.Conn) {
	// Send request to server
	var (
		DepartmentName string
	)

	fmt.Print("Enter department name: ")
	fmt.Scanln(&DepartmentName)

	item := Department{id: 1, DepartmentName: DepartmentName}
	createDepartment(conn, item)
}
func createSendToServerStud(conn net.Conn) {
	// Send request to server
	var (
		StudentName  string
		DepartmentID int
	)

	fmt.Print("Enter student name: ")
	fmt.Scanln(&StudentName)

	fmt.Print("Enter Department ID: ")
	fmt.Scanln(&DepartmentID)

	item := Student{id: 1, StudentName: StudentName, DepartmentID: DepartmentID}
	createStudent(conn, item)
}
func createDepartment(conn net.Conn, department Department) {
	request := []byte("insertDepartment ")

	x := []byte(strconv.Itoa(department.id))
	x = append(x, " "...)
	request = append(request, x...)

	x = []byte(department.DepartmentName)
	x = append(x, " "...)
	request = append(request, x...)

	conn.Write(request)
}
func createStudent(conn net.Conn, student Student) {
	request := []byte("insertStudent ")

	x := []byte(student.StudentName)
	x = append(x, " "...)
	request = append(request, x...)

	x = []byte(strconv.Itoa(student.DepartmentID))
	x = append(x, " "...)
	request = append(request, x...)

	conn.Write(request)
}
func getAllDepartments(conn net.Conn) {
	request := []byte("queryAllDepartments")
	conn.Write(request)
	get := make([]byte, 1024)
	conn.Read(get)
	fmt.Println(string(get))
}
func getAllStudents(conn net.Conn) {
	request := []byte("queryAllStudents")
	conn.Write(request)
	get := make([]byte, 1024)
	conn.Read(get)
	fmt.Println(string(get))
}

func deleteDepartment(conn net.Conn) {
	var ID int
	fmt.Print("Enter id: ")
	fmt.Scanln(&ID)

	request := []byte("deleteDepartment ")
	x := []byte(strconv.Itoa(ID))
	request = append(request, x...)

	conn.Write(request)
	get := make([]byte, 1024)
	conn.Read(get)
	fmt.Println(string(get))
}
func deleteStudent(conn net.Conn) {

	var ID int
	fmt.Print("Enter id: ")
	fmt.Scanln(&ID)

	request := []byte("deleteStudent ")
	x := []byte(strconv.Itoa(ID))
	request = append(request, x...)

	conn.Write(request)
	get := make([]byte, 1024)
	conn.Read(get)
	fmt.Println(string(get))
}
func updateDepartment(conn net.Conn) {

	var (
		DepartmentName string
		id             int
	)

	fmt.Print("Enter id: ")
	fmt.Scanln(&id)
	fmt.Print("Enter department name: ")
	fmt.Scanln(&DepartmentName)

	department := Department{id: id, DepartmentName: DepartmentName}

	request := []byte("updateDepartment ")

	x := []byte(strconv.Itoa(department.id))
	x = append(x, " "...)
	request = append(request, x...)

	x = []byte(department.DepartmentName)
	x = append(x, " "...)
	request = append(request, x...)

	conn.Write(request)

	get := make([]byte, 1024)
	conn.Read(get)
	fmt.Println(string(get))
}
func updateStudent(conn net.Conn) {

	var (
		StudentName  string
		DepartmentID int
		id           int
	)

	fmt.Print("Enter id: ")
	fmt.Scanln(&id)
	fmt.Print("Enter Student name: ")
	fmt.Scanln(&StudentName)
	fmt.Print("Enter Department ID: ")
	fmt.Scanln(&DepartmentID)

	student := Student{id: id, StudentName: StudentName, DepartmentID: DepartmentID}

	request := []byte("updateStudent ")

	x := []byte(strconv.Itoa(student.id))
	x = append(x, " "...)
	request = append(request, x...)

	x = []byte(student.StudentName)
	x = append(x, " "...)
	request = append(request, x...)

	x = []byte(strconv.Itoa(student.DepartmentID))
	x = append(x, " "...)
	request = append(request, x...)

	conn.Write(request)

	get := make([]byte, 1024)
	conn.Read(get)
	fmt.Println(string(get))
}
