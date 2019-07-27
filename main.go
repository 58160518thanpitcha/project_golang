package main

import (
    "encoding/json"
    "log"
    "net/http"
    "github.com/gorilla/mux"
    "time"
	"github.com/yugabyte/gocql"
	"fmt"
	"strconv"
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
  // Connect to the cluster.
  cluster := gocql.NewCluster("127.0.0.1", "127.0.0.2", "127.0.0.3")

  // Use the same timeout as the Java driver.
  cluster.Timeout = 12 * time.Second

  // Create the session.
  session, _ := cluster.CreateSession()
  defer session.Close()

  // Set up the keyspace and table.
  if err := session.Query("CREATE KEYSPACE IF NOT EXISTS ybtest").Exec(); err != nil {
	  log.Fatal(err)
  }
  fmt.Println("Created keyspace ybtest")


  if err := session.Query(`DROP TABLE IF EXISTS ybtest.number`).Exec(); err != nil {
	  log.Fatal(err)
  }
  var createStmt = `CREATE TABLE ybtest.number (id int PRIMARY KEY, 
														 firstnum varchar, 
														 lastnum varchar,
														 result varchar)`;
  if err := session.Query(createStmt).Exec(); err != nil {
	  log.Fatal(err)
  }
  fmt.Println("Created table ybtest.number")

    params := mux.Vars(req)
	var numeric Numeric
	json.NewDecoder(req.Body).Decode(&numeric) //แปลงข้อมูลjsonที่รับมาเป็น struct 
	// Insert into the table.
	num1 := strconv.Itoa(numeric.Firstnum)
	num2 := strconv.Itoa(numeric.Lastnum)
	
	var insertStmt string = "INSERT INTO ybtest.number(id,firstnum, lastnum)" + 
		" VALUES(1,'"+num1+"', '"+num2+"')";
	if err := session.Query(insertStmt).Exec(); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted data: %s\n", insertStmt)

	 // Read from the table.
	 var firstnum string
	 var lastnum string
	 iter := session.Query(`SELECT firstnum, lastnum FROM ybtest.number WHERE id = 1`).Iter()//query ข้อมูลออกมา
	 fmt.Printf("Query for id=1 returned: ");
	 for iter.Scan(&firstnum, &lastnum) {
		 fmt.Printf("Row[%s, %s]\n", firstnum, lastnum)
	 }
	 
	 if err := iter.Close(); err != nil {
		 log.Fatal(err)
	 }
	 n1, _ := strconv.Atoi(firstnum)
	 n2, _ := strconv.Atoi(lastnum)
	 
	var result int = n1+n2; //คำนวณผลบวก
	resultdb := strconv.Itoa(result)
	var update string = "UPDATE ybtest.number SET result = '"+resultdb+"' WHERE id = 1"; //เพิ่มผลลัพธ์ที่่คำนวณได้เข้า db
	if err := session.Query(update).Exec(); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted Result: %s\n", update)
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
