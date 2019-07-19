package main
import (
    "encoding/json"
    "fmt"
    "net/http"
)
type number struct {
    Numfirst int
    Numlast  int
}
func getNumberAll(w http.ResponseWriter, r *http.Request) {
    addNumber := number{
                Numfirst: 2,
                Numlast:  4,
              }
    json.NewEncoder(w).Encode(addNumber)
}
func homePage(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Welcome to the HomePage!")
}
func handleRequest() {
    http.HandleFunc("/", homePage)
    http.HandleFunc("/getAddress", getNumberAll)
    http.ListenAndServe(":8080", nil)
}
func main() {
    handleRequest()
}