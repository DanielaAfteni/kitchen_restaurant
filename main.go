package main

import (
	// Package json implements encoding and decoding of JSON.
	// The mapping between JSON and Go values is described in the documentation for the Marshal and Unmarshal functions.

	// Package ioutil implements some I/O utility functions.

	// Package os provides a platform-independent interface to operating system functionality.

	// importing the gin, because is a high-performance HTTP web framework written in Golang (Go).

	"bytes"
	"encoding/json"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

const TIME_UNIT = 250

func main() {
	router := gin.Default()
	router.POST("/order", getOrder)

	rand.Seed(time.Now().UnixNano())
	AllCooks = append(AllCooks, &Cook{
		Id:          1,
		Rank:        3,
		Proficiency: 4,
		Name:        "Gordon Ramsay",
		CatchPhrase: "Hey, panini head, are you listening to me?",
	})
	AllCooks = append(AllCooks, &Cook{
		Id:          2,
		Rank:        2,
		Proficiency: 3,
		Name:        "Igor Dodon",
		CatchPhrase: "Oameni buni!",
	})
	AllCooks = append(AllCooks, &Cook{
		Id:          3,
		Rank:        2,
		Proficiency: 2,
		Name:        "Maia Sandu",
		CatchPhrase: "Scapam de coruptie.",
	})
	AllCooks = append(AllCooks, &Cook{
		Id:          4,
		Rank:        1,
		Proficiency: 2,
		Name:        "LaLa Land",
		CatchPhrase: "Mergem la mare!!!",
	})
	for cook := range AllCooks {
		for i := 0; i < AllCooks[cook].Proficiency; i++ {
			go AllCooks[cook].AreCooking()
		}
	}
	router.Run(":8081")
}

func getOrder(c *gin.Context) {
	var order *Order
	if err := c.BindJSON(&order); err != nil {
		log.Err(err).Msg("Error!!")
		return
	}
	log.Printf("The new order came to the kitchen %+v", order)
	FoodsListWhichAreReady = append(FoodsListWhichAreReady, ReadyFoods{
		orderId:   order.OrderId,
		orderSize: 0,
	})
	FoodsDivision(order)
	Order_list = append(Order_list, order)
	c.IndentedJSON(http.StatusCreated, order)
}

func gatherBackAllFinishedFoods(food EachFoodFromOrder, cookId int) {
	for eFLWAR := range FoodsListWhichAreReady {
		if food.orderId == FoodsListWhichAreReady[eFLWAR].orderId {

			FoodsListWhichAreReady[eFLWAR].orderSize++
			FoodsListWhichAreReady[eFLWAR].foods = append(FoodsListWhichAreReady[eFLWAR].foods, FoodCookedByCook{
				FoodId: food.id,
				CookId: cookId,
			})

			if FoodsListWhichAreReady[eFLWAR].orderSize == food.orderSize {
				for eachOrderID := range Order_list {
					if Order_list[eachOrderID].OrderId == food.orderId {
						orderPrepared := &OrderPrepared{
							Order:          *Order_list[eachOrderID],
							CookingTime:    time.Now().Unix() - Order_list[eachOrderID].PickUpTime,
							CookingDetails: FoodsListWhichAreReady[eachOrderID].foods,
						}
						log.Printf("Finally prepared order is %+v", orderPrepared)
						Order_list = Order_list[1:]
						jsonBody, err := json.Marshal(orderPrepared)
						if err != nil {
							log.Err(err).Msg("Error!!!")
						}
						contentType := "application/json"
						//_, err = http.Post("http://dining_hall_restaurant:8080/distribution", contentType, bytes.NewReader(jsonBody))
						_, err = http.Post("http://localhost:8080/distribution", contentType, bytes.NewReader(jsonBody))
						if err != nil {
							log.Err(err).Msg("Error!!")
							return
						}
						break
					}
				}
			}
			break
		}
	}
}
