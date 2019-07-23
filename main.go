package main

import (
    "encoding/json"
    "log"
    "net/http"
    "github.com/gorilla/mux"
)

type Numeric struct { //สร้าง struct
    ID        string   `json:"id,omitempty"`
    Firstnum  int   `json:"firstnum,omitempty"`
	Lastnum   int   `json:"lastnum,omitempty"`
	Result    int    `json:"result,omitempty"`
    
}

var number []Numeric

func GetAllNumberEndpoint(w http.ResponseWriter, req *http.Request) {//ฟังก์ชันแสดงข้อมูลทั้งหมดที่addไว้
    params := mux.Vars(req)
    for _, item := range number {
        if item.ID == params["id"] {
            json.NewEncoder(w).Encode(item)
            return
        }
    }
    json.NewEncoder(w).Encode(&Numeric{})
}

func GetNumberEndpoint(w http.ResponseWriter, req *http.Request) { //ฟังก์ชันแสดงข้อมูลตาม ID ที่addไว้
    json.NewEncoder(w).Encode(number)
}

func CreateNumberEndpoint(w http.ResponseWriter, req *http.Request) {
    params := mux.Vars(req)
	var numeric Numeric
    json.NewDecoder(req.Body).Decode(&numeric) //แปลงข้อมูลjsonที่รับมาเป็น struct
    var result int = numeric.Firstnum+numeric.Lastnum; //คำนวณผลบวก
    numeric.Result = result; //เก็บผลบวกใส่ไว้ใน struct
    numeric.ID = params["id"]
	number = append(number, numeric)
    json.NewEncoder(w).Encode(number)// แปลง struct เป็น json เพื่อแสดงผล
}
func DeleteNumberEndpoint(w http.ResponseWriter, req *http.Request) { //ฟังก์ชันDeleteข้อมูล
    params := mux.Vars(req)
    for index, item := range number {
        if item.ID == params["id"] {
            number = append(number[:index], number[index+1:]...)
            break
        }
    }
    json.NewEncoder(w).Encode(number)
}

func main() {
    router := mux.NewRouter()
    router.HandleFunc("/list", GetAllNumberEndpoint).Methods("GET")
    router.HandleFunc("/list/{id}", GetNumberEndpoint).Methods("GET")
    router.HandleFunc("/list/{id}", CreateNumberEndpoint).Methods("POST")
    router.HandleFunc("/list/{id}", DeleteNumberEndpoint).Methods("DELETE")
    log.Fatal(http.ListenAndServe(":12345", router))
}