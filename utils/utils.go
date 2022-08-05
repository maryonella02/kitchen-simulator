package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"kitchen-simulator/models"
	"log"
	"os"
)

func ReadCooks(filename string) models.Cooks {
	jsonFile, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}
	log.Printf("Successfully Opened %s \n", filename)

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var cooks models.Cooks

	err = json.Unmarshal(byteValue, &cooks)
	if err != nil {
		log.Panic("Cannot unmarshal json cooks")
	}

	return cooks
}

func ReadMapOfMenus(filename string) map[int]models.Dish {

	jsonFile, err := os.Open(filename)
	if err != nil {
		log.Println(err)
	}
	log.Printf("Successfully Opened %s \n", filename)

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var dishes models.Dishes

	err = json.Unmarshal(byteValue, &dishes)
	if err != nil {
		log.Panic("Cannot unmarshal json menu")
	}

	return ConvertListOfDishesToMap(dishes)
}

func ConvertListOfDishesToMap(dishes models.Dishes) map[int]models.Dish {
	mapOfDishes := make(map[int]models.Dish)
	for _, dish := range dishes.Dishes {
		mapOfDishes[dish.ID] = dish
	}

	return mapOfDishes
}

func RemoveFromSliceByIndex(slice []models.Order, s int) []models.Order {
	return append(slice[:s], slice[s+1:]...)
}
