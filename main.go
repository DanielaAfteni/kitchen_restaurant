package main

import (
	// Package json implements encoding and decoding of JSON.
	// The mapping between JSON and Go values is described in the documentation for the Marshal and Unmarshal functions.
	"encoding/json"
	// Package ioutil implements some I/O utility functions.
	"io/ioutil"
	// Package os provides a platform-independent interface to operating system functionality.
	"os"
	// importing the gin, because is a high-performance HTTP web framework written in Golang (Go).
	"github.com/DanielaAfteni/kitchen_restaurant/basicproces"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// setting the configuration
	cfg := config()
	basicproces.SettingtheConfig(cfg)
	// Channels are a typed conduit through which you can send and receive values with the channel operator, <-.
	// we create an order channel
	orderChan := make(chan basicproces.Order, cfg.NrOfTables)
	//we create the menu, by calling the corresponding function
	menu := basicproces.GetMenu()
	orderList := basicproces.NewOrderList(orderChan, menu)
	orderList.Run()

	r := gin.Default()
	// The Kitchen should handle requests of receiving orders from the Dinning Hall
	// and add received order to order list .
	r.POST("/order", func(c *gin.Context) {
		var order basicproces.Order
		if err := c.ShouldBindJSON(&order); err != nil {
			log.Err(err).Msg("Error binding JSON")
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		orderChan <- order
		c.JSON(200, gin.H{"message": "Order received"})
	})
	r.Run(":8081")
}

func config() basicproces.Config {
	// Output writes the output for a logging event
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Logger = log.With().Caller().Logger()
	// To open and read the json file
	file, err := os.Open("config/scfg.json")
	// in case of an error, it returns a message
	if err != nil {
		log.Fatal().Err(err).Msg("Error opening menu.json")
	}
	// close the file
	defer file.Close()

	byteValue, _ := ioutil.ReadAll(file)
	var scfg basicproces.Config
	json.Unmarshal(byteValue, &scfg)

	return scfg
}
