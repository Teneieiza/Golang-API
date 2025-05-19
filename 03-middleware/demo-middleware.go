package middleware

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

func findID(Id int) (*Pet, int) {
	for i, pet := range PetList {
		if pet.Id == Id {
			return &pet, i
		}

	}
	return nil, 0
}

func petHandler(w http.ResponseWriter, r *http.Request) {
	urlPathSegment := strings.Split(r.URL.Path, "pet/")
	Id, err := strconv.Atoi(urlPathSegment[len(urlPathSegment)-1])

	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	pet, listItemIndex := findID(Id)

	if pet == nil {
		http.Error(w, fmt.Sprintf("no pet with this Id: %d", Id), http.StatusNotFound)
		return
	}

	switch r.Method {

	case http.MethodGet:
		petJson, err := json.Marshal(pet)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(petJson)

	case http.MethodPut:
		var updatedPet Pet
		byteBody, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = json.Unmarshal(byteBody, &updatedPet)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if updatedPet.Id != Id {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		pet = &updatedPet
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

//สร้าง func สำหรับทำ middleware ให้รับค่ามาเป็น Handler และ return ออกมาเป็น handler
func middlewareHandler(handler http.Handler) http.Handler{
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		fmt.Println("before handler middleware start")
		
		//------------------------------------------

		// ข้อมูลที่เราต้องการให้ middleware ทำ

		//------------------------------------------

		handler.ServeHTTP(w, r)
		fmt.Println("middleware finished")
	})
}

func Middleware() {
	//สร้างต้วแปรสำหรับเก็บ handler
	petListHandler := http.HandlerFunc(petsHandler)
	petItemHandler := http.HandlerFunc(petHandler)

	//ในส่วนของ path ก็ให้ใน midleware มาครอบไว้
	//เมื่อลองยิงตาม path ก็จะเห็น midleware log ออกมา
	http.Handle("/pet", middlewareHandler(petListHandler))
	http.Handle("/pet/", middlewareHandler(petItemHandler))
	http.ListenAndServe(":5000", nil)
}
