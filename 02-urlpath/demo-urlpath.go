package urlpath

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
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

func getNextID() int {
	highestId := -1
	for _, pet := range PetList {
		if highestId < pet.Id {
			highestId = pet.Id
		}
	}
	return highestId + 1
}

// Func สำหรับหา Id
// ซึ่งจะวนลูปใน PetList และ ซึ่งถ้าเจอ Id จาก Pet ที่ตรงกับ Id ที่ป้อนเข้ามา
// จะ return pet และ Index ออกไป
// แต่ถ้าไม่เจอจะ return ค่าว่างและ 0
func findID(Id int) (*Pet, int) {
	for i, pet := range PetList {
		if pet.Id == Id {
			return &pet, i
		}

	}
	return nil, 0
}

// Func Handler สำหรับเรียกใช้งาน Method CRUD
// ในที่นี้ petHandler จะเป็น method สำหรับ path ที่เรียกตาม Id
func petHandler(w http.ResponseWriter, r *http.Request) {
	//ทำการ split เพื่อแบ่ง path โดยแบ่งจาก pet/ เช่น localhost:5000/pet/6
	//ก็จะแบ่งได้เป็น localhost:5000/ และ 6
	urlPathSegment := strings.Split(r.URL.Path, "pet/")
	//หลังจากนั้นสร้างตัวแปร Id มาเก็บข้อมูลที่ได้จากการแบ่งมา -1 ก็หมายถึงตัวหลังสุด
	// เช่น เช่น localhost:5000/pet/6 แบ่งมาจะได้ localhost:5000/ และ 6 เก็บตัวหลังสุดก็คือ 6 ไว้ใน Id
	Id, err := strconv.Atoi(urlPathSegment[len(urlPathSegment)-1])
	//สร้างเงื่อนไขดักขึ้นมา ถ้ากรณีที่มี error ให้เข้าเงื่อนไข
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	//ซึ่ง func findID จะ return ค่าออกมา 2 ตัวทำให้ต้องสร้างตัวแปรมารับค่า 2 ตัว
	pet, listItemIndex := findID(Id)
	//กรณีที่ pet ไม่ได้รับค่ามาหรือ nil ก็ให้เข้าเงื่อนไข error
	if pet == nil {
		http.Error(w, fmt.Sprintf("no pet with this Id: %d", Id), http.StatusNotFound)
		return
	}
	//สร้าง Swithcase สำหรับ method ต่างๆ
	switch r.Method {
	//สร้าง method GET
	case http.MethodGet:
		//แปลงค่า pet ที่รับมาเป็น json
		petJson, err := json.Marshal(pet)
		//กรณีเกิด error ก็ให้แสดง error
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//ถ้าไม่เจอ error ก็ให้ request ส่ง header ไป และให้ respone แสดงข้อมูล json
		w.Header().Set("Content-Type", "application/json")
		w.Write(petJson)

		//สร้าง method PUT
	case http.MethodPut:
		//เริ่มจากการสร้างตัวแปรสำหรับเก็บข้อมูลที่จะ update
		var updatedPet Pet
		// หลังจากนั้นให้ทำการอ่านข้อมูลใน Body และ return มาเก็บไว้ให้ byteBody
		// กรณีเจอ error ก็ให้เป็นไว้ในตัวแปร err
		byteBody, err := io.ReadAll(r.Body)
		//ถ้าพบ error ก็ให้เข้าเงื่อนไขแสดง StatusBadRequest
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		//แปลงค่าจาก Body ที่รับมา จาก json ให้เป็น object
		//กรณีที่เกิด error ก็ให้เก็บค่าไว้ใน err
		err = json.Unmarshal(byteBody, &updatedPet)
		//ถ้าพบ error ก็ให้เข้าเงื่อนไขแสดง StatusBadRequest
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		//updatedPet.Id ไม่เท่ากับ Id ให้แสดง error ออกมา
		//StatusBadRequest
		if updatedPet.Id != Id {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		//กรณีไม่เจอ error ก็ให้ข้อมูลจาก updatedPet ไปใส่ใน pet
		pet = &updatedPet
		//โดยการเข้าถึงข้อมูลของ PetList[] โดยใช้ index จาก listItemIndex
		//ที่ดึงผ่าน findId
		PetList[listItemIndex] = *pet
		w.WriteHeader(http.StatusOK)
		return

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}


func petsHandler(w http.ResponseWriter, r *http.Request) {
	petJson, err := json.Marshal(PetList)
	switch r.Method {

	case http.MethodGet:
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(petJson)
		return

	case http.MethodPost:
		var newPet Pet
		bodybytes, err := io.ReadAll(r.Body)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = json.Unmarshal(bodybytes, &newPet)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if newPet.Id != 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		newPet.Id = getNextID()
		PetList = append(PetList, newPet)
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Pet added successfully!"))
		return
	}
}

func UrlPath() {
	//ประกาศ path และเรียกใช้งาน handler(Function การทำงานของ method ต่างๆ)
	// handlerFunc สำหรับ path ที่เรียกทั้งหมด
	http.HandleFunc("/pet", petsHandler)
	// handlerFunc สำหรับ path ที่เรียกเฉพาะรายการ
	http.HandleFunc("/pet/", petHandler)
	http.ListenAndServe(":5000", nil)
}
