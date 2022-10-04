package main

type Order struct {
	TableId            int     `json:"table_id"`
	OrderId            int     `json:"order_id"`
	WaiterId           int     `json:"waiter_id"`
	Priority           int     `json:"priority"`
	MenuItemIds        []int   `json:"items"`
	MaxPreparationTime float32 `json:"max_wait"`
	PickUpTime         int64   `json:"pick_up_time"`
}

type OrderPrepared struct {
	Order
	CookingTime    int64              `json:"cooking_time"`
	CookingDetails []FoodCookedByCook `json:"cooking_details"`
}

var Order_list []*Order
var OrderPreparedList []*Order
