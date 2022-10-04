package main

import (
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

//const NrCooks int = 3

const NrCooks int = 4

var AllCooks []*Cook

type Cook struct {
	Id          int
	Rank        int
	Proficiency int
	Name        string
	CatchPhrase string
}
type Rank int

var CookingApparatus map[string]int
var semStove = make(chan int, 1)
var semOven = make(chan int, 2)

var OrderMutex sync.Mutex

// process for cooking apparatus
func cookingApparatusProcess(foodList []EachFoodFromOrder) int {
	if CookingApparatus["oven"] == 0 && CookingApparatus["stove"] == 0 {
		if len(foodList) == 1 && foodList[0].Food.cookingApparatus != "" {
			if foodList[0].Food.cookingApparatus == "oven" {
				semOven <- 1
				CookingApparatus["oven"] -= 1
			} else if foodList[0].Food.cookingApparatus == "stove" {
				semStove <- 1
				CookingApparatus["stove"] -= 1
			}
			return 0
		}
		for foodID := range foodList {
			if foodList[foodID].Food.cookingApparatus == "" {
				return foodID
			}
		}
	} else {
		for foodID := range foodList {
			if foodList[foodID].Food.cookingApparatus == "" {
				return foodID
			} else if CookingApparatus[foodList[foodID].Food.cookingApparatus] != 0 {
				if foodList[foodID].Food.cookingApparatus == "oven" {
					semOven <- 1
					CookingApparatus["oven"] -= 1
				} else if foodList[foodID].Food.cookingApparatus == "stove" {
					semStove <- 1
					CookingApparatus["stove"] -= 1
				}
				return foodID
			}
		}
	}

	return 0
}

func getOrderListItem(rank int) EachFoodFromOrder {
	var eachFood EachFoodFromOrder
	if rank == 3 && len(FoodList3.GetFoodList()) > 0 {
		idOfEachFood := cookingApparatusProcess(FoodList3.GetFoodList())
		eachFood = FoodList3.GetFoodList()[idOfEachFood]
		FoodList3.ReduceFoodList(idOfEachFood)
	} else if rank >= 2 && len(FoodList2.GetFoodList()) > 0 {
		idOfEachFood := cookingApparatusProcess(FoodList2.GetFoodList())
		eachFood = FoodList2.GetFoodList()[idOfEachFood]
		FoodList2.ReduceFoodList(idOfEachFood)
	} else {
		idOfEachFood := cookingApparatusProcess(FoodList1.GetFoodList())
		eachFood = FoodList1.GetFoodList()[idOfEachFood]
		FoodList1.ReduceFoodList(idOfEachFood)
	}
	FoodToPrepare--
	return eachFood
}

func (theCook *Cook) AreCooking() {
	CookingApparatus = map[string]int{
		"oven":  2,
		"stove": 1,
	}
	for {
		OrderMutex.Lock()
		if (theCook.Rank == 1 && len(FoodList1.GetFoodList()) < 1) || (theCook.Rank == 2 && len(FoodList2.GetFoodList()) < 1) {
			OrderMutex.Unlock()
			continue
		}
		if FoodToPrepare > 0 {
			food := getOrderListItem(theCook.Rank)
			OrderMutex.Unlock()
			log.Printf("START to cook by cook %d the food %+v\n", theCook.Id, food)
			<-time.After(time.Duration(food.preparationTime) * time.Second)
			CookingApparatus[food.Food.cookingApparatus] += 1
			if food.Food.cookingApparatus == "oven" {
				<-semOven
			} else if food.Food.cookingApparatus == "stove" {
				<-semStove
			}
			log.Printf("FINISH to cook by cook %d the food %+v\n", theCook.Id, food)
			gatherBackAllFinishedFoods(food, theCook.Id)
		} else {
			OrderMutex.Unlock()
		}
	}
}
