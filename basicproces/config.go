// Packages are used to organize related Go source files together into a single unit
// The current go program is in package basicproces
package basicproces

// A structure or struct in Golang is a user-defined type that allows to group/combine items of possibly different types into a single type.
// Any real-world entity which has some set of properties/fields can be represented as a struct.
// This concept is generally compared with the classes in object-oriented programming.

// Config - has the majority of the properties/fields as int (because are integer numbers)
// and which are taken from the json file
type Config struct {
	TimeUnit      int    `json:"time_unit"`
	NrOfTables    int    `json:"nr_of_tables"`
	DiningHallUrl string `json:"dining_hall_url"`
}

// set for this configuration including all the variables of the application (2 integers, 1 url)
var scfg Config = Config{
	TimeUnit:      1000,
	NrOfTables:    10,
	DiningHallUrl: "http://dining_hall_restaurant:8080",
}

// funtion for setting the configuration
func SettingtheConfig(s Config) {
	scfg = s
}
