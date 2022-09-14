// Packages are used to organize related Go source files together into a single unit
// The current go program is in package basicproces
package basicproces

import (
	// Package json implements encoding and decoding of JSON.
	// The mapping between JSON and Go values is described in the documentation for the Marshal and Unmarshal functions.
	"encoding/json"
	// Package ioutil implements some I/O utility functions.
	"io/ioutil"
	"os"

	"github.com/rs/zerolog/log"
)

// A structure or struct in Golang is a user-defined type that allows to group/combine items of possibly different types into a single type.
// Any real-world entity which has some set of properties/fields can be represented as a struct.
// This concept is generally compared with the classes in object-oriented programming.

// Menu - has just 2 properties/fields as int and array/slice of food
// and which are taken from the json file
type Menu struct {
	// This Menu struct will hold all food and count variables of the application that we read from file
	FoodsCount int
	Foods      []Food
}

// Food - has the majority of the properties/fields as int (because are integer numbers)
// and which are taken from the json file
type Food struct {
	// This Food struct will hold all food variables of the application that we read from file
	Id               int    `json:"id"`
	Name             string `json:"name"`
	PreparationTime  int    `json:"preparation_time"`
	Complexity       int    `json:"complexity"`
	CookingApparatus string `json:"cooking_apparatus"`
}

// function for creation a menu
func GetMenu() Menu {
	// To open and read the json file
	file, err := os.Open("config/menu.json")
	// in case of an error, it returns a message
	if err != nil {
		log.Fatal().Err(err).Msg("Error appeared at opening menu.json. Try to find it.")
	}
	// regarding a file, that we opened, we need to close it when we’re done. Here’s how we could do that with defer.
	defer file.Close()
	//ReadAll - reads from file and returns the data it read.
	byteValue, _ := ioutil.ReadAll(file)
	// we create a variable of type Menu
	var menu Menu
	// we use the Unmarshal() function in package encoding/json - to unpack or decode the data from JSON to a struct (menu).
	json.Unmarshal(byteValue, &menu)
	// we set the lenght of the menu - as the lenght of all the foods that we have
	menu.FoodsCount = len(menu.Foods)
	// we return the menu
	return menu
}
