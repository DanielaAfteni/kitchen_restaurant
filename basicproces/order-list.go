package basicproces

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/rs/zerolog/log"
)

type OrderList struct {
	Distributions     map[int]*Distribution
	ReceiveOrder      <-chan Order
	ReceiveCookedFood chan FoodOrder
	Cooks             []*Cook
	Menu              Menu
}

type CooksDetails struct {
	Cooks []CookDetails
}

func NewOrderList(receiveOrder <-chan Order, menu Menu) *OrderList {
	ol := &OrderList{
		Distributions:     make(map[int]*Distribution),
		ReceiveOrder:      receiveOrder,
		ReceiveCookedFood: make(chan FoodOrder),
		Menu:              menu,
	}
	// To open and read the json file
	file, err := os.Open("config/cooks.json")
	// in case of an error, it returns a message
	if err != nil {
		log.Fatal().Err(err).Msg("Error occured while opening cooks.json. Try to find it.")
	}
	//ReadAll - reads from file and returns the data it read.
	byteValue, _ := ioutil.ReadAll(file)
	// we create a variable of type CooksDetails
	var cooksDetails CooksDetails
	// we use the Unmarshal() function in package encoding/json - to unpack or decode the data from JSON to a struct (cooksDetails).
	json.Unmarshal(byteValue, &cooksDetails)

	ol.Cooks = make([]*Cook, len(cooksDetails.Cooks))
	// we are looping over the indexes and elements in details about the cooks (especially its name)
	for i, cookDetails := range cooksDetails.Cooks {
		ol.Cooks[i] = NewCook(i, cookDetails, ol.ReceiveCookedFood, ol.Menu)
		// Gordon Ramsay entered the kitchen cook_id=0
		// <specific cook> entered the kitchen <its specific id>
		log.Info().Int("cook_id", i).Msgf("%s entered the kitchen", cookDetails.Name)
	}

	return ol
}

// function to run order list
func (ol *OrderList) Run() {
	// have 2 calls of food being cooked
	go ol.SendFoodOrderToCooks()
	// and food coming back from the cook
	go ol.ReceiveFoodOrderFromCooks()
}

func (ol *OrderList) SendFoodOrderToCooks() {
	// we are looping over elements in received order
	for order := range ol.ReceiveOrder {
		// we take the values and variables which were distributed
		ol.Distributions[order.OrderId] = &Distribution{
			Order:          order,
			CookingTime:    time.Now().UnixMilli(),
			CookingDetails: make([]CookingDetail, 0),
			ReceivedItems:  make([]bool, len(order.Items)),
		}
		// show that the kitchen received order
		log.Info().Int("order_id", order.OrderId).Msg("The kitchen received order")
		// we are looping over elements in order items
		for i, id := range order.Items {
			// take the food
			food := ol.Menu.Foods[id-1]
			// and set the Sent food as false
			IsFoodOrderSent := false
			// untill the sent food is true, we are doing:
			for !IsFoodOrderSent {
				// we are looping over elements in order list cooks
				for _, cook := range ol.Cooks {
					// in case if the cook can cook
					if cook.CanCook(food) {
						// take the corresponding food
						foodOrder := FoodOrder{
							OrderId: order.OrderId,
							ItemId:  i,
							CookingDetail: CookingDetail{
								FoodId: food.Id,
								CookId: cook.Id,
							},
						}
						// cook is going to be occupied by the cooking process
						atomic.AddInt64(&cook.Occupation, 1)
						// start cooking
						go cook.CookFood(foodOrder)
						// show that food is assigned to a cook
						log.Info().Int("order_id", order.OrderId).Int("item_id", i).Int("food_id", food.Id).Int("cook_id", cook.Id).Msgf("%s order assigned to %s", food.Name, cook.Name)
						// set that food is done
						IsFoodOrderSent = true
						break
					}
				}
			}
		}
	}
}

// function for food coming back from the cook
func (ol *OrderList) ReceiveFoodOrderFromCooks() {
	// we are looping over elements in cooked order by cook
	for foodOrder := range ol.ReceiveCookedFood {
		// set a distribution for the food to the order
		distribution := ol.Distributions[foodOrder.OrderId]
		// in case if the food is wrong
		if distribution.Order.Items[foodOrder.ItemId] != foodOrder.FoodId {
			// then we show that
			log.Warn().Int("order_id", foodOrder.OrderId).Int("item_id", foodOrder.ItemId).Msg("There is a received wrong food item")
			continue
		}
		// in case if the food is already received
		if distribution.ReceivedItems[foodOrder.ItemId] {
			// then we show that
			log.Warn().Int("order_id", foodOrder.OrderId).Int("item_id", foodOrder.ItemId).Msg("There is a food item already received")
			continue
		}
		// set the food - received
		distribution.ReceivedItems[foodOrder.ItemId] = true
		// we add to the distribution cooking details the food order cooking details
		distribution.CookingDetails = append(distribution.CookingDetails, foodOrder.CookingDetail)
		//if the length of the CookingDetails are the same as the nr of items
		if len(distribution.CookingDetails) == len(distribution.Order.Items) {
			distribution.CookingTime = (time.Now().UnixMilli() - distribution.CookingTime) / int64(scfg.TimeUnit)
			jsonBody, err := json.Marshal(distribution)
			if err != nil {
				log.Fatal().Err(err).Msg("Error marshalling distribution")
			}
			contentType := "application/json"
			_, err = http.Post(scfg.DiningHallUrl+"/distribution", contentType, bytes.NewReader(jsonBody))
			if err != nil {
				log.Fatal().Err(err).Msg("Error sending distribution to dining hall")
			}
			log.Info().Int("order_id", foodOrder.OrderId).Msg("Distribution sent to dining hall")
			delete(ol.Distributions, foodOrder.OrderId)
		}
	}
}
