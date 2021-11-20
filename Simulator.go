package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Order struct {
	ID         int   `json:"id"`
	Items      []int `json:"items"`
	TableID    int   `json:"table_id"`
	WaiterID   int   `json:"waiter_id"`
	Priority   int   `json:"priority"`
	MaxWait    int   `json:"max_wait"`
	PickUpTime int64 `json:"pick_up_time"`
}
type ReadyOrder struct {
	ID          int   `json:"id"`
	Items       []int `json:"items"`
	TableID     int   `json:"table_id"`
	WaiterID    int   `json:"waiter_id"`
	Priority    int   `json:"priority"`
	MaxWait     int   `json:"max_wait"`
	PickUpTime  int64 `json:"pick_up_time"`
	CookingTime int64 `json:"cooking_time"`
}

func main() {
	go func() {
		for {
			go func() {
				url := "http://localhost:8081/test"
				req, err := http.NewRequest("GET", url, nil)
				client := &http.Client{}
				resp, err := client.Do(req)
				if err != nil {
					print("skip the error")
				} else {
					defer resp.Body.Close()
					fmt.Println("Request sent")
				}

			}()
			time.Sleep(time.Second)
		}
	}()

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
	worker(order)
}

func worker(order Order) {
	request := getJsonRequest(order)
	makeRequest(request)

}

func getJsonRequest(order Order) []byte {
	readyOrder := &ReadyOrder{ID: order.ID,
		Items:       order.Items,
		Priority:    order.Priority,
		MaxWait:     order.MaxWait,
		WaiterID:    order.WaiterID,
		TableID:     order.TableID,
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
