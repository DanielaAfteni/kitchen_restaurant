// Packages are used to organize related Go source files together into a single unit
// The current go program is in package basicproces
package basicproces

import (
	// Package atomic provides low-level atomic memory primitives useful for implementing synchronization algorithms.
	"sync/atomic"
	// Package time provides functionality for measuring and displaying time.
	"time"

	"github.com/rs/zerolog/log"
)

// A structure or struct in Golang is a user-defined type that allows to group/combine items of possibly different types into a single type.
// Any real-world entity which has some set of properties/fields can be represented as a struct.
// This concept is generally compared with the classes in object-oriented programming.

// FoodOrder - has just 2 properties/fields as int and Cooking Detail
type FoodOrder struct {
	CookingDetail
	ItemId  int
	OrderId int
}

// CookDetails - has just 2 properties/fields as int and 2 strings of name details
// and which are taken from the json file
type CookDetails struct {
	Rank        int    `json:"rank"`
	Proficiency int64  `json:"proficiency"`
	Name        string `json:"name"`
	CatchPhrase string `json:"catch_phrase"`
}

// Menu - has just 2 properties/fields as int, the menu and a channel for finally prepared food
type Cook struct {
	CookDetails
	Id             int
	Occupation     int64
	SendCookedFood chan<- FoodOrder
	Menu           Menu
}

// function to return all the cooks with their information, data
func NewCook(id int, cookDetails CookDetails, sendChan chan<- FoodOrder, menu Menu) *Cook {
	return &Cook{
		CookDetails:    cookDetails,
		Id:             id,
		Occupation:     0,
		SendCookedFood: sendChan,
		Menu:           menu,
	}
}

// boolean function for determination if the cook can cook
func (c *Cook) CanCook(food Food) bool {
	// if Proficiency bigger than Occupation
	isFree := atomic.LoadInt64(&c.Occupation) < c.Proficiency
	// if Rank bigger or equal than Complexity of the food
	isQualified := food.Complexity <= c.Rank
	return isFree && isQualified
}

// function for cooks to cook
func (c *Cook) CookFood(foodOrder FoodOrder) {
	// select the food from food oredr from menu
	food := c.Menu.Foods[foodOrder.FoodId-1]
	// if Complexity bigger than Rank
	if food.Complexity > c.Rank {
		//it means that cooker is not qualified
		log.Warn().Int("cook_id", c.Id).Msgf("%s is not qualified to cook %s", c.Name, food.Name)
		atomic.AddInt64(&c.Occupation, -1)
		return
	}
	// if Occupation bigger than Proficiency
	if atomic.LoadInt64(&c.Occupation) > c.Proficiency {
		// it means that the cook is too busy
		log.Warn().Int("cook_id", c.Id).Msgf("%s is too busy", c.Name)
		atomic.AddInt64(&c.Occupation, -1)
		return
	}
	// calculate preparation time
	preparationTime := time.Duration(food.PreparationTime * scfg.TimeUnit * int(time.Millisecond))
	// sleep a bit at preparation time
	time.Sleep(preparationTime)
	// show that the food was cooked by the cook
	log.Info().Int("cook_id", c.Id).Int("food_id", foodOrder.FoodId).Int("order_id", foodOrder.OrderId).Msgf("%s cooked %s", c.Name, food.Name)
	// send data to the finally prepared food
	c.SendCookedFood <- foodOrder
	atomic.AddInt64(&c.Occupation, -1)
}
