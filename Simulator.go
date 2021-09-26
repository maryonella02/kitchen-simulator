package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type Order struct {
	ID         int   `json:"id"`
	Items      []int `json:"items"`
	Priority   int   `json:"priority"`
	MaxWait    int   `json:"max_wait"`
	PickUpTime int64 `json:"pick_up_time"`
}
type ReadyOrder struct {
	ID          int   `json:"id"`
	Items       []int `json:"items"`
	Priority    int   `json:"priority"`
	MaxWait     int   `json:"max_wait"`
	PickUpTime  int64 `json:"pick_up_time"`
	CookingTime int64 `json:"cooking_time"`
}

func main() {
	http.HandleFunc("/order", HandleRequest)
	log.Fatal(http.ListenAndServe(":8082", nil))
}

func HandleRequest(rw http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var order Order
	err := decoder.Decode(&order)
	if err != nil {
		panic(err)
	}
	log.Println(order)
	fmt.Println("Request Handled")
	var wg sync.WaitGroup
	for i := 1; i <= 1000; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			worker(order)
		}()
	}
	wg.Wait()
}

func worker(order Order) {
	request := getJsonRequest(order)
	time.Sleep(time.Second)
	makeRequest(request)

}

func getJsonRequest(order Order) []byte {
	readyOrder := &ReadyOrder{ID: order.ID,
		Items:       order.Items,
		Priority:    order.Priority,
		MaxWait:     order.MaxWait,
		PickUpTime:  order.PickUpTime,
		CookingTime: getCookingTime(order)}
	b, err := json.Marshal(readyOrder)
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
	return b
}

func getCookingTime(order Order) int64 {
	cookingTime := getUnixTimestamp() - order.PickUpTime
	return cookingTime
}

func getUnixTimestamp() int64 {
	now := time.Now()
	sec := now.Unix()
	return sec
}
func makeRequest(b []byte) {
	url := "http://localhost:8081/distribution"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Println(string(b))
	fmt.Println("Request sent")

}
