package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"kitchen-simulator/models"
	"kitchen-simulator/utils"
	"net/http"
	"time"
)

var Orders *models.OrdersList
var Cooks []models.Cook
var ApparatusList models.CookingApparatusQueue
var DishesMenu map[int]models.Dish

func convertOrderToReadyOrder(order models.Order, cookingTime int) []byte {
	readyOrder := &models.ReadyOrder{ID: order.ID,
		Items:       order.Items,
		Priority:    order.Priority,
		MaxWait:     order.MaxWait,
		WaiterID:    order.WaiterID,
		TableID:     order.TableID,
		PickUpTime:  order.PickUpTime,
		CookingTime: cookingTime,
	}
	b, err := json.Marshal(readyOrder)
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
	return b
}

func StartServingDishes(cook models.Cook) {
	fmt.Println("Received name ", cook.Name)
	for {
		time.Sleep(time.Second)
		fmt.Println("Before lock")
		Orders.Lock()
		fmt.Println("As cook ", cook.Name, "I try to find some orders")
		if len(Orders.AllOrders) > 0 {
			//TODO: Implement finding food by his rank
			for idx, order := range Orders.AllOrders {
				receiptList, canCook := CanCookOrder(cook, order)
				if canCook {
					Orders.AllOrders = utils.RemoveFromSliceByIndex(Orders.AllOrders, idx)
					fmt.Println("Unlocked")
					Orders.Unlock()
					StartPrepareOrder(cook, order, receiptList)
					break
				} else {
					Orders.Unlock()
				}
			}
		} else {
			fmt.Println("There are no orders availalbe")
			Orders.Unlock()
		}

	}

}

func CanCookOrder(cook models.Cook, order models.Order) ([]models.Dish, bool) {
	receiptList := makeReceiptList(order)
	canCook := true
	fmt.Println("Finding order can cook")
	for _, food := range receiptList {
		if cook.Rank < food.Complexity {
			canCook = false
			break
		}
	}
	fmt.Println("Can cook order id ", order.ID, canCook)
	return receiptList, canCook

}

func makeReceiptList(order models.Order) []models.Dish {
	foodReceipts := make([]models.Dish, len(order.Items))
	for idx, foodID := range order.Items {
		receipt := DishesMenu[foodID]
		foodReceipts[idx] = receipt
	}

	return foodReceipts
}

func StartPrepareOrder(cook models.Cook, order models.Order, foodReceipts []models.Dish) {
	fmt.Println("Starting prepare order ", order.ID)
	startPreparationTime := time.Now()

	for _, food := range foodReceipts {
		fmt.Println("Start prepared food id with name", food.ID, food.Name)
		switch food.CookingApparatus {
		case models.STOVE:
			prepareUsingFreeApparat(ApparatusList.Stoves, food.PreparationTime)
		case models.OVEN:
			prepareUsingFreeApparat(ApparatusList.Ovens, food.PreparationTime)
		default:
			time.Sleep(time.Duration(food.PreparationTime) * time.Millisecond)
		}
	}
	elapsed := time.Since(startPreparationTime).Milliseconds()
	fmt.Println("Done ", order.ID)
	sendPreparedOrder(order, int(elapsed))
}

func prepareUsingFreeApparat(apparats chan models.Apparat, preparationTime int) {
	fmt.Println("Before using apparat")
	freeApparat := <-apparats
	fmt.Println("Apparat tooked")
	time.Sleep(time.Duration(preparationTime) * time.Millisecond)
	apparats <- freeApparat
	fmt.Println("aparat released")

}

func sendPreparedOrder(order models.Order, cookedTime int) {
	fmt.Println("Sending order ", order.ID)
	bytes := convertOrderToReadyOrder(order, cookedTime)
	makeRequest(bytes)
	fmt.Println("Sended order ", order.ID)
	// fmt.Println("Deleting order from lis", order.ID)
	// Orders.Lock()
	// utils.RemoveFromSliceByOrderID(Orders.AllOrders, order.ID)
	// Orders.Unlock()

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
