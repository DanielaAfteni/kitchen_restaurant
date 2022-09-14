// The current go program is in package basicproces
package basicproces

// Distribution - has just one of the properties/fields as int (because are integer numbers)
// and which are taken from the json file
type Distribution struct {
	Order
	CookingTime    int64           `json:"cooking_time"`
	CookingDetails []CookingDetail `json:"cooking_details"`
	ReceivedItems  []bool          `json:"-"`
}

// CookingDetail - has just one of the properties/fields as int (because are integer numbers)
// and which are taken from the json file
type CookingDetail struct {
	FoodId int `json:"food_id"`
	CookId int `json:"cook_id"`
}
