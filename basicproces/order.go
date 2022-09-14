// Packages are used to organize related Go source files together into a single unit
// The current go program is in package basicproces
package basicproces

// Order - has the majority of the properties/fields as int (because are integer numbers)
// and which are taken from the json file
type Order struct {
	// This Order struct will hold all order details variables of the application that we read from file
	OrderId    int     `json:"order_id"`
	TableId    int     `json:"table_id"`
	WaiterId   int     `json:"waiter_id"`
	Items      []int   `json:"items"`
	Priority   int     `json:"priority"`
	MaxWait    float64 `json:"max_wait"`
	PickUpTime int64   `json:"pick_up_time"`
}
