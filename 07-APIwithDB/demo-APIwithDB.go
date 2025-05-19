package apiwithdb

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var database *sql.DB

var basePath = "/api"
var petFoodPath = "petfood"

type PetFood struct {
	Id         int     `json:"id"`
	Name       string  `json:"name"`
	Price      float64 `json:"price"`
	AnimalType string  `json:"animaltype"`
	FoodType   string  `json:"foodtype"`
	ImageUrl   string  `json:"imageurl"`
}

// ----------------------------------Function SetupDB--------------------------------
// ฟังก์ชันนี้ใช้สำหรับตั้งค่าการเชื่อมต่อกับฐานข้อมูล MySQL
func SetupDB() {
	var err error
	// ประกาศตัวแปร database เพื่อเก็บการเชื่อมต่อกับฐานข้อมูล
	// โดยใช้ sql.Open() เพื่อเปิดการเชื่อมต่อกับฐานข้อมูล MySQL
	database, err = sql.Open("mysql", "root:Ten0826817189@tcp(127.0.0.1:3306)/petdatabase")

	if err != nil {
		panic(err.Error())
	}

	fmt.Println("DataBase : ", database)
	// SetConnMaxLifetime ตั้งค่าช่วงเวลาที่การเชื่อมต่อจะถูกเก็บไว้ใน pool
	// โดยจะทำการปิดการเชื่อมต่อที่ไม่ได้ใช้งานหลังจากเวลาที่กำหนด
	database.SetConnMaxLifetime(time.Minute * 3)

	// SetMaxOpenConns ตั้งค่าจำนวนการเชื่อมต่อสูงสุดที่สามารถเปิดได้พร้อมกัน
	// โดยจะทำการปิดการเชื่อมต่อที่เกินจำนวนที่กำหนด
	database.SetMaxOpenConns(10)

	// SetMaxIdleConns ตั้งค่าจำนวนการเชื่อมต่อที่ idle (ไม่ได้ใช้งาน) ที่จะถูกเก็บไว้ใน pool
	// โดยจะทำการปิดการเชื่อมต่อที่ idle ที่เกินจำนวนที่กำหนด
	database.SetMaxIdleConns(10)

	fmt.Println("Connected to the database successfully!")
}

//-----------------------------------------------------------------------------------

// ---------------------------------Function getAllFoods-----------------------------
// ฟังก์ชันนี้ใช้สำหรับดึงข้อมูลทั้งหมดจากตาราง petfood
func getAllFoods() ([]PetFood, error) {
	//ประกาศ context เพื่อใช้ในการกำหนดเวลา timeout สำหรับการ query ข้อมูล
	// โดยจะทำการยกเลิกการ query หากใช้เวลานานเกิน 3 วินาที
	context, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// result เก็บผลลัพธ์จากการ query ข้อมูลจากตาราง petfood
	// โดยใช้คำสั่ง SELECT เพื่อดึงข้อมูล Id, Name, Price, AnimalType, FoodType, ImageUrl
	// กรณีที่ error ให้ทำการปิดการทำงานและแสดง error
	result, err := database.QueryContext(context, "SELECT Id, Name, Price, AnimalType, FoodType, ImageUrl FROM petfood")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer result.Close()

	// ประกาศตัวแปร petFoods เพื่อเก็บข้อมูลที่ query ได้
	// โดยใช้ make() เพื่อสร้าง slice ที่มีขนาดเริ่มต้นเป็น 0
	// และใช้ append() เพื่อเพิ่มข้อมูลที่ query ได้ลงใน slice
	petFoods := make([]PetFood, 0)
	for result.Next() {
		var petFood PetFood
		if err := result.Scan(&petFood.Id, &petFood.Name, &petFood.Price, &petFood.AnimalType, &petFood.FoodType, &petFood.ImageUrl); err != nil {
			log.Fatal(err)
			return nil, err
		}
		petFoods = append(petFoods, petFood)
	}
	return petFoods, nil
}

//-----------------------------------------------------------------------------------

// ---------------------------------Function insertPetFood---------------------------
// ฟังก์ชันนี้ใช้สำหรับเพิ่มข้อมูลใหม่ลงในตาราง petfood
// โดยใช้คำสั่ง INSERT INTO เพื่อเพิ่มข้อมูล Id, Name, Price, AnimalType, FoodType, ImageUrl
// กรณีที่ error ให้ทำการปิดการทำงานและแสดง error
func insertPetFood(petFood PetFood) (int64, error) {
	// ประกาศ context เพื่อใช้ในการกำหนดเวลา timeout สำหรับการ query ข้อมูล
	// โดยจะทำการยกเลิกการ query หากใช้เวลานานเกิน 3 วินาที
	context, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// result เก็บผลลัพธ์จากการ query ข้อมูลจากตาราง petfood
	// โดยใช้คำสั่ง INSERT INTO เพื่อเพิ่มข้อมูล Id, Name, Price, AnimalType, FoodType, ImageUrl
	result, err := database.ExecContext(context, "INSERT INTO petfood (Name, Price, AnimalType, FoodType, ImageUrl) VALUES (?, ?, ?, ?, ?)",
		petFood.Name,
		petFood.Price,
		petFood.AnimalType,
		petFood.FoodType,
		petFood.ImageUrl)

	if err != nil {
		log.Println("Insert error:", err)
		return 0, err
	}

	insertId, err := result.LastInsertId()
	if err != nil {
		log.Println("LastInsertId error:", err)
		return 0, err
	}

	return insertId, nil
}

//----------------------------------------------------------------------------------

// ---------------------------------Function getPetFoodById---------------------------
// ฟังก์ชันนี้ใช้สำหรับดึงข้อมูลจากตาราง petfood โดยใช้ Id เป็นตัวระบุ
func getPetFoodById(Id int) (*PetFood, error) {
	context, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	row := database.QueryRowContext(context, "SELECT Id, Name, Price, AnimalType, FoodType, ImageUrl FROM petfood WHERE Id = ?", Id)
	PetFood := &PetFood{}
	err := row.Scan(&PetFood.Id, &PetFood.Name, &PetFood.Price, &PetFood.AnimalType, &PetFood.FoodType, &PetFood.ImageUrl)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return PetFood, nil
}

//-----------------------------------------------------------------------------------

// ---------------------------------Function removePetFoodById-----------------------
// ฟังก์ชันนี้ใช้สำหรับลบข้อมูลจากตาราง petfood โดยใช้ Id เป็นตัวระบุ
func removePetFoodById(Id int) error {
	context, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := database.ExecContext(context, "DELETE FROM petfood WHERE Id = ?", Id)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

// -------------------------------------Handle-PetFoods---[many]---------------------
func handlePetFoods(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	//---------------------------------GET--------------------------------------------
	case http.MethodGet:
		petFoodList, err := getAllFoods()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, "Error retrieving pet foods:", err)
			return
		}

		petfoodsJson, err := json.Marshal(petFoodList)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, "Error marshalling JSON:", err)
			return
		}

		_, err = w.Write(petfoodsJson)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, "Error writing response:", err)
			return
		}
	//-------------------------------------------------------------------------------

	//-------------------------------POST--------------------------------------------
	case http.MethodPost:
		var petFood PetFood

		err := json.NewDecoder(r.Body).Decode(&petFood)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "Error decoding JSON:", err)
			return
		}

		Id, err := insertPetFood(petFood)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, "Error inserting pet food:", err)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write(fmt.Appendf(nil, "Pet food created with ID: %d", Id))
		return
	//-------------------------------------------------------------------------------

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintln(w, "Method not allowed")
		return
	}
}

//-----------------------------------------------------------------------------------

// ---------------------------------Handle-Petfood---[one]--------------------------
func handlePetfood(w http.ResponseWriter, r *http.Request) {
	urlPathSegments := strings.Split(r.URL.Path, fmt.Sprintf("/%s/", petFoodPath))
	if len(urlPathSegments[1:]) > 1 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	Id, err := strconv.Atoi(urlPathSegments[len(urlPathSegments)-1])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	switch r.Method {

	//---------------------------------GET--------------------------------------------
	case http.MethodGet:
		petfood, err := getPetFoodById(Id)
		if err != nil {
			http.Error(w, "Pet food not found", http.StatusNotFound)
			return
		}

		if petfood == nil {
			http.Error(w, "Pet food not found", http.StatusNotFound)
			return
		}

		petfoodJson, err := json.Marshal(petfood)
		if err != nil {
			http.Error(w, "Error marshalling JSON", http.StatusInternalServerError)
			return
		}

		_, err = w.Write(petfoodJson)
		if err != nil {
			http.Error(w, "Error writing response", http.StatusInternalServerError)
			return
		}
	//-------------------------------------------------------------------------------

	//-------------------------------DELETE--------------------------------------------
	case http.MethodDelete:
		err := removePetFoodById(Id)
		if err != nil {
			http.Error(w, "Error deleting pet food", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(fmt.Appendf(nil, "Pet food %d was Deleted", Id))
		return
	//-------------------------------------------------------------------------------

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

//-----------------------------------------------------------------------------------

// -------------------------------CORS Middleware-----------------------------------
// ฟังก์ชันนี้ใช้สำหรับตั้งค่า CORS (Cross-Origin Resource Sharing) เพื่ออนุญาตให้มีการเข้าถึง API จากโดเมนอื่นๆ
func corsMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Authorization, X-Requested-With")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		handler.ServeHTTP(w, r)
	})
}

//-----------------------------------------------------------------------------------

// -------------------------------SetupRoutes----------------------------------------
// ฟังก์ชันนี้ใช้สำหรับตั้งค่าเส้นทาง (routes) ของ API โดยใช้ http.Handle() เพื่อกำหนดเส้นทางและ handler ที่จะทำงานเมื่อมีการเรียกใช้งานเส้นทางนั้นๆ
func SetupRoutes(apiBasePath string) {
	// ตั้งค่าเส้นทางสำหรับ petfoods
	petfoodsHandler := http.HandlerFunc(handlePetFoods)
	http.Handle(fmt.Sprintf("%s/%s", apiBasePath, petFoodPath), corsMiddleware(petfoodsHandler))

	// ตั้งค่าเส้นทางสำหรับ petfood โดยใช้ Id เป็นตัวระบุ
	petfoodHandler := http.HandlerFunc(handlePetfood)
	http.Handle(fmt.Sprintf("%s/%s/", apiBasePath, petFoodPath), corsMiddleware(petfoodHandler))
}

//-----------------------------------------------------------------------------------

func ApiWithDB() {
	SetupDB()
	SetupRoutes(basePath)
	log.Fatal(http.ListenAndServe(":5000", nil))
}
