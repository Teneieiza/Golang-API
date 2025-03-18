package webservice

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type Pet struct {
	Id     int     `json:"id"`
	Name   string  `json:"name"`
	Weight float64 `json:"weight"`
	Type   string  `json:"type"`
}

var PetList []Pet

func init() {
	PetJson := `[
		{
			"id":1,
			"name":"Tigger",
			"weight":2.7,
			"type":"Cat"
		},
		{
			"id":2,
			"name":"Toufu",
			"weight":2.4,
			"type":"Cat"
		},
		{
			"id":3,
			"name":"Chinchan",
			"weight":0.7,
			"type":"GuiniePig"
		},
		{
			"id":4,
			"name":"LittleDragon",
			"weight":0.5,
			"type":"Hamster"
		},
		{
			"id":5,
			"name":"Yellow",
			"weight":0.3,
			"type":"Hamster"
		}		
	]`
	err := json.Unmarshal([]byte(PetJson), &PetList)
	if err != nil {
		log.Fatal(err)
	}
}

// Func สำหรับการ auto id
// โดยการกำหนดค่าเริ่มต้นเท่ากับ -1 เนื่องด้วย id ไม่มีทางเป็นค่าติดลบ
// จากนั้นนำมา loop for โดยกำหนดการวนอยู่ใน PetList
// ซึ่งถ้า highestId มีค่าน้อยกว่า pet.Id ก็จะให้ highestId มีค่าเท่ากับ pet.Id
// และจะวนจนกว่าค่า highestId มีค่าไม่น้อยกว่า pet.Id
// จากนั้นจะให้ highestId + 1
func getNextID() int {
	highestId := -1
	for _, pet := range PetList {
		if highestId < pet.Id {
			highestId = pet.Id
		}
	}
	return highestId + 1
}

// Func Handler สำหรับเรียกใช้งาน Method CRUD
func petHandler(w http.ResponseWriter, r *http.Request) {
	//แปลง PetList ให้อยู่ในรูปแบบ Json
	petJson, err := json.Marshal(PetList)
	switch r.Method {

	//method GET
	case http.MethodGet:
		//กรณีเกิด error ก็จะโยน error กลับมาให้
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		//กรณีที่ไม่เจอ error ก็จะแสดง "Content-Type", "application/json" กลับมา
		w.Header().Set("Content-Type", "application/json")
		//ใช้ Write เพื่อแสดงข้อมูล Json ออกมา
		w.Write(petJson)

	//method POST
	case http.MethodPost:
		//ประกาศตัวแปรสำหรับเพิ่มข้อมูล
		var newPet Pet
		//สร้างตัวแปล bodybytes มารับค่าที่อ่านผ่าน request ที่อยู่ใน body ทั้งหมด
		bodybytes, err := io.ReadAll(r.Body)

		//กรณีเกิด error ก็จะโยน error กลับมาให้
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		//แปลงค่า Json ที่รับมาจาก body เป็น struct และเพิ่มเข้าไปใน newPet
		err = json.Unmarshal(bodybytes, &newPet)
		//กรณีเกิด error ก็จะโยน error กลับมาให้
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		//กรณีเกิด error ก็จะโยน error กลับมาให้
		if newPet.Id != 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		//ถ้าไม่ error ก็จะทำการเพิ่ม newPet เข้าไปไว้ใน PetList
		//ซึ่งจะมีการกำหนด auto id ผ่าน func getNextID
		newPet.Id = getNextID()
		PetList = append(PetList, newPet)
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Pet added successfully!"))
		return
	}
}

func WorkRequest() {
	//ประกาศ path และเรียกใช้งาน handler(Function การทำงานของ method ต่างๆ)
	http.HandleFunc("/pet", petHandler)
	//ประกาศ port ที่จะใช้งาน
	http.ListenAndServe(":5000", nil)
}
