package gowithmysql

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

//---------------------------Fuction createTable---------------------------
// ฟังก์ชันนี้ใช้สำหรับสร้างตาราง members ในฐานข้อมูล petdatabase
// โดยมีคอลัมน์ Id, Name, PetId, Username, Password และ CreatedAt
// โดยมี Id เป็น PRIMARY KEY และ PetId เป็น FOREIGN KEY ที่อ้างอิงจากตาราง pet
func createMemberTable(db *sql.DB) {
	query := `CREATE TABLE members (
		Id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
		Name TEXT NOT NULL,
		PetId INT NOT NULL,
		Username TEXT NOT NULL,
		Password TEXT NOT NULL,
		CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (PetId) REFERENCES pet(Id)
	);`

	if _, err := db.Exec(query); err != nil {
		fmt.Println("Error creating table:", err)
		return
	}
}
//--------------------------------------------------------------------------


//---------------------------Fuction insertMember---------------------------
// ฟังก์ชันนี้ใช้สำหรับเพิ่มข้อมูลสมาชิกใหม่ลงในตาราง members
func insertMember(db *sql.DB) {
	var (
		name     string
		username string
		password string
		petId    int
	)

	fmt.Print("Enter your name: ")
	fmt.Scanln(&name)

	fmt.Print("Enter your username: ")
	fmt.Scanln(&username)

	fmt.Print("Enter your password: ")
	fmt.Scanln(&password)

	fmt.Print("Enter your pet ID: ")
	if _, err := fmt.Scanf("%d", &petId); err != nil {
		fmt.Println("Invalid Pet ID")
		os.Exit(1)
	}

	// โดยมีการตรวจสอบว่า petId ที่ป้อนเข้ามามีอยู่ในตาราง pet หรือไม่
	// ถ้ามีจะทำการเพิ่มข้อมูลสมาชิกใหม่ลงในตาราง members
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM pet WHERE Id = ?)", petId).Scan(&exists)
	if err != nil || !exists {
		fmt.Println("Pet ID not found in pet table.")
		return
	}

	// สั่ง insert โดยใช้ SELECT จาก pet
	query := `
	INSERT INTO members (Name, PetId, Username, Password, CreatedAt)
	SELECT ?, ?, ?, ?, ? FROM pet WHERE Id = ?
	`

	_, err = db.Exec(query, name, petId, username, password, time.Now(), petId)
	if err != nil {
		fmt.Println("Error inserting member:", err)
		return
	}

	fmt.Println("Member inserted successfully!")
}
//--------------------------------------------------------------------------



//---------------------------Fuction deleteMember---------------------------
// ฟังก์ชันนี้ใช้สำหรับลบข้อมูลสมาชิกจากตาราง members โดยใช้ Id เป็นเงื่อนไขในการลบ
// โดยจะทำการลบข้อมูลสมาชิกที่มี Id ตรงกับที่ผู้ใช้ป้อนเข้ามา
func deleteMember(db *sql.DB) {
	var deleteId int
	fmt.Print("Enter the ID of the member to delete: ")
	fmt.Scanln(&deleteId)
	_, err := db.Exec("DELETE FROM members WHERE Id = ?", deleteId)
	if err != nil {
		fmt.Println("Error deleting member:", err)
		return
	}
	fmt.Println("Member deleted successfully!")
}
//--------------------------------------------------------------------------


//---------------------------Fuction query by Id---------------------------
// ฟังก์ชันนี้ใช้สำหรับ query ข้อมูลจากตาราง pet โดยใช้ Id เป็นเงื่อนไขในการค้นหา
// โดยจะทำการแสดงผลข้อมูลที่ได้จากการ query ออกมา
func queryPet(db *sql.DB) {
	var (
		inputId int
		Id      int
		Name    string
		Weight  float64
		Type    string
	)

	fmt.Print("Please enter the pet ID you want to query: ")
	if _, err := fmt.Scanf("%d", &inputId); err != nil {
		fmt.Println("Invalid input. Please enter a valid number:", err)
		os.Exit(1)
	}

	query := "SELECT Id, Name, Weight, Type FROM pet WHERE Id = ?"
	err := db.QueryRow(query, inputId).Scan(&Id, &Name, &Weight, &Type)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No pet found with the given ID.")
			return
		}
		fmt.Println("Error querying database:", err)
		return
	}

	fmt.Printf("Id: %d, Name: %s, Weight: %.2f, Type: %s\n", Id, Name, Weight, Type)
}
//--------------------------------------------------------------------------


func Gowithmysql() {
	database, err := sql.Open("mysql", "root:Ten0826817189@tcp(127.0.0.1:3306)/petdatabase")
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		os.Exit(1)
	}

	defer database.Close()
	fmt.Println("Connected to database successfully")

	// fmt.Println("Database-Not query:", database)
	// queryPet(database)
	// createMemberTable(database)
	// insertMember(database)
	// deleteMember(database)
}
