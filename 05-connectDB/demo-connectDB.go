package connectDB

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
)

	// Query the database
	// โดยการประกาศตัวแปร Id, Name, Weight, Type เพื่อเก็บค่าที่ได้จากการ query
func queryDatabase(db *sql.DB, id int) {
	var (
		Id     int
		Name   string
		Weight float64
		Type   string
	)

	// สร้าง query string เพื่อดึงข้อมูลจากตาราง pet โดยใช้ Id เป็นเงื่อนไขในการค้นหา
	// โดยใช้ db.QueryRow เพื่อ query ข้อมูลเพียง 1 แถว
	query := "SELECT Id, Name, Weight, Type FROM pet WHERE Id = ?"
	if err := db.QueryRow(query, id).Scan(&Id, &Name, &Weight, &Type); err != nil {
		fmt.Println("Error querying database:", err)
		return
	}
	fmt.Printf("Id: %d, Name: %s, Weight: %.2f, Type: %s\n", Id, Name, Weight, Type)
}

func ConnectDB() {
	//สร้างตัวแปร inputId เพื่อเก็บค่าที่ผู้ใช้ป้อนเข้ามา
	// โดยใช้ fmt.Scanf เพื่ออ่านค่าจาก standard input
	var inputId int

	fmt.Println("Please enter the pet ID you want to query:")
	if _, err := fmt.Scanf("%d", &inputId); err != nil {
		fmt.Println("Invalid Input Please input number:", err)
		os.Exit(1)
	}

	// เชื่อมต่อกับฐานข้อมูล MySQL
	// โดยใช้ sql.Open เพื่อเปิดการเชื่อมต่อกับฐานข้อมูล
	databse, err := sql.Open("mysql", "root:Ten0826817189@tcp(127.0.0.1:3306)/petdatabase")
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		os.Exit(1)
	}
	
	//ปิดการเชื่อมต่อหลังใช้งานเสร็จ
	defer databse.Close()
	fmt.Println("Connected to database successfully")

	// fmt.Println("Database-Not query:", databse)
	queryDatabase(databse, inputId)
}
