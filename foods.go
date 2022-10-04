package main

import (
	"sort"
	"sync"

	"github.com/rs/zerolog/log"
)

type Food struct {
	id               int
	name             string
	preparationTime  float32
	complexity       int
	cookingApparatus string
}

type EachFoodFromOrder struct {
	orderId       int
	orderSize     int
	orderPriority int
	Food
}

type FoodLists struct {
	foodListMutex sync.Mutex
	foodList      []EachFoodFromOrder
}

type FoodCookedByCook struct {
	FoodId int `json:"food_id"`
	CookId int `json:"cook_id"`
}

type ReadyFoods struct {
	orderId   int
	orderSize int
	foods     []FoodCookedByCook
}

var FoodsListWhichAreReady []ReadyFoods
var FoodList1 FoodLists
var FoodList2 FoodLists
var FoodList3 FoodLists
var FoodToPrepare = 0

func (f *FoodLists) GetLockUnlockFoodList() {
	f.foodListMutex.Lock()
	defer f.foodListMutex.Unlock()
}

func (f *FoodLists) GetFoodList() []EachFoodFromOrder {
	f.GetLockUnlockFoodList()
	return f.foodList
}

func (f *FoodLists) SetFoodList(order *Order, id int) {
	f.GetLockUnlockFoodList()
	// like to show
	f.foodList = append(f.foodList, EachFoodFromOrder{
		orderId:       order.OrderId,
		orderSize:     len(order.MenuItemIds),
		orderPriority: order.Priority,
		Food:          Menu[id],
	})
	f.PriorityFixing(order.Priority)
}

func (f *FoodLists) ReduceFoodList(id int) {
	f.GetLockUnlockFoodList()
	f.foodList = append(f.foodList[:id], f.foodList[id+1:]...)
}

func minimalOrder(foodList []EachFoodFromOrder) int {
	min := 5
	for _, food := range foodList {
		min = determineTheLowestPriority(food.orderPriority, min)
	}
	return min
}

func determineTheLowestPriority(currentPrio int, lowestPriority int) int {
	if currentPrio < lowestPriority {
		lowestPriority = currentPrio
	}
	return lowestPriority
}

func (f *FoodLists) PriorityFixing(existingOrderPriority int) {
	sort.SliceStable(f.foodList, func(p, q int) bool {
		return f.foodList[p].orderPriority > f.foodList[q].orderPriority
	})
	if (existingOrderPriority - minimalOrder(f.foodList)) > 2 {
		for _, food := range f.foodList {
			if food.orderPriority != 5 {
				food.orderPriority += 1
			}
		}
	}
}

// by their complexity
func FoodsDivision(order *Order) {
	for _, eachItemID := range order.MenuItemIds {
		if Menu[eachItemID-1].complexity == 1 {
			FoodList1.SetFoodList(order, eachItemID-1)
		} else if Menu[eachItemID-1].complexity == 2 {
			FoodList2.SetFoodList(order, eachItemID-1)
		} else if Menu[eachItemID-1].complexity == 3 {
			FoodList3.SetFoodList(order, eachItemID-1)
		} else {
			log.Printf("Wrong complexity")
		}
		FoodToPrepare++
	}
}

var Menu = []Food{
	{
		id:               1,
		name:             "pizza",
		preparationTime:  20,
		complexity:       2,
		cookingApparatus: "oven",
	},
	{
		id:               2,
		name:             "salad",
		preparationTime:  10,
		complexity:       1,
		cookingApparatus: "",
	},
	{
		id:               3,
		name:             "zeama",
		preparationTime:  7,
		complexity:       1,
		cookingApparatus: "stove",
	},
	{
		id:               4,
		name:             "Scallop Sashimi with Meyer Lemon Confit",
		preparationTime:  32,
		complexity:       3,
		cookingApparatus: "",
	},
	{
		id:               5,
		name:             "Island Duck with Mulberry Mustard",
		preparationTime:  35,
		complexity:       3,
		cookingApparatus: "oven",
	},
	{
		id:               6,
		name:             "Waffles",
		preparationTime:  10,
		complexity:       1,
		cookingApparatus: "stove",
	},
	{
		id:               7,
		name:             "Aubergine",
		preparationTime:  20,
		complexity:       2,
		cookingApparatus: "oven",
	},
	{
		id:               8,
		name:             "Lasagna",
		preparationTime:  30,
		complexity:       2,
		cookingApparatus: "oven",
	},
	{
		id:               9,
		name:             "Burger",
		preparationTime:  15,
		complexity:       1,
		cookingApparatus: "stove",
	},
	{
		id:               10,
		name:             "Gyros",
		preparationTime:  15,
		complexity:       1,
		cookingApparatus: "",
	},
}
