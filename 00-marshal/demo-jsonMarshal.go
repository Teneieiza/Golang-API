package marshal

import (
	"encoding/json"
	"fmt"
)

type petMarshal struct {
	Id        int
	PetName   string
	PetWeight float64
	TypeOfPet string
}

func Marshal() {
	//สร้างตัวแปรเก็บข้อมูลในรูปแบบ Struct
	myAnimalStruct := petMarshal{101, "Tigger", 2.7, "Cat"}

	//นำเข้าไปทำงานใน Marshal เพื่อแปลง Struct ไปเป็น Object,Json
	data, _ := json.Marshal(&myAnimalStruct)
	
	fmt.Println(string(data))
}
