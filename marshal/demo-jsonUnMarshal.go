package marshal

import (
	"encoding/json"
	"fmt"
	"log"
)

type petUnmarshal struct {
	Id        int
	PetName   string
	PetWeight float64
	TypeOfPet string
}

func UnMarshal() {
	//ประกาศตัวแปรเพื่อเก็บชนิดของข้อมูล
	petHome := petUnmarshal{}

	//สร้างตัวแปรเก็บข้อมูลของ Object,Json
	myAnimalJson := []byte(`{"Id":101,"PetName":"Tigger","PetWeight":2.7,"TypeOfPet":"Cat"}`)
	
	//นำเข้าไปทำงานใน Unmarshal เพื่อแปลง Object,Json ไปเป็น Struct
	err := json.Unmarshal(myAnimalJson, &petHome)
	if err != nil {
		log.Fatal(err)
	}
	//การเข้าถึงข้อมูลทั้งหมด
	fmt.Println(petHome)

	//การเข้าถึงข้อมูลใน object
	fmt.Println(petHome.PetName)
}
