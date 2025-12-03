package delivery

import "fmt"

func DeliveryTimeFactory(tt string) DeliveryTimeStrategy {
	switch tt {
	case "on_foot":
		return FootTimeStrategy{}
	case "scooter":
		return ScooterTimeStrategy{}
	case "car":
		return CarTimeStrategy{}
	default:
		fmt.Printf("the type of transport (%s) was not found", tt)
		return FootTimeStrategy{}
	}
}
