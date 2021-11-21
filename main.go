package main

import (
	"encoding/json"
	"fmt"
	"kitchen-simulator/models"
	"kitchen-simulator/services"
	"kitchen-simulator/utils"
	"log"
	"net/http"
	"sort"
)

const (
	Ovens  = 2
	Stoves = 1
)

func main() {
	services.DishesMenu = utils.ReadMapOfMenus("menu.json")
	services.Orders = &models.OrdersList{
		AllOrders: make([]models.Order, 0),
	}

	initCookingApparatus()

	services.Cooks = utils.ReadCooks("cooks.json").Cooks

	for _, cooker := range services.Cooks {
		fmt.Println("Cooker is", cooker)
	}

	fmt.Println("There are ", len(services.Cooks), "cooks")
	for _, cook := range services.Cooks {
		fmt.Println("Starting as ", cook.Name)
		go services.StartServingDishes(cook)
	}

	http.HandleFunc("/order", HandleRequest)
	log.Fatal(http.ListenAndServe(":8082", nil))
}

func HandleRequest(rw http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var order models.Order
	err := decoder.Decode(&order)
	if err != nil {
		panic(err)
	}
	log.Println(order)
	fmt.Println("Request Handled locked")
	// Add order in order list
	services.Orders.Lock()
	fmt.Println("Before appending new order")
	services.Orders.AllOrders = append(services.Orders.AllOrders, order)
	fmt.Println("Before sorting")
	sort.Sort(models.ByPriorityAndMaxWait(services.Orders.AllOrders))
	fmt.Println("Now len of total order list is ", len(services.Orders.AllOrders))
	fmt.Println("Unlocked request handled")
	defer services.Orders.Unlock()

}

func initCookingApparatus() {
	services.ApparatusList = models.CookingApparatusQueue{
		Ovens:  make(chan models.Apparat, Ovens),
		Stoves: make(chan models.Apparat, Stoves),
	}

	for i := 0; i < Ovens; i++ {
		services.ApparatusList.Ovens <- models.Apparat{}
	}

	for i := 0; i < Stoves; i++ {
		services.ApparatusList.Stoves <- models.Apparat{}
	}

}
