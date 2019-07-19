package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	
)
const input = `
{
	"msg": {
		"numfirst": "175",
		"numlast": "221"
	}
}
`

type Envelope struct {
	Type string
	Msg  interface{}
}
func main() {
	var env Envelope
	if err := json.Unmarshal([]byte(input), &env); err != nil {
		log.Fatal(err)
	}
	// for the love of Gopher DO NOT DO THIS
	var num1 string = env.Msg.(map[string]interface{})["numfirst"].(string)
	var num2 string = env.Msg.(map[string]interface{})["numlast"].(string)

	a, err := strconv.Atoi(num1)
	if err == nil {
		fmt.Println("a: ",a)
	}
	b, err := strconv.Atoi(num2)
	if err == nil {
		fmt.Println("b: ",b)
	}
	var c int = a+b;
	fmt.Println("result: a+b = ",c)

}